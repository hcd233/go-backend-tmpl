package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/hcd233/aris-api-tmpl/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// LogSamplingRule 日志采样规则
//
//	@author centonhuang
//	@update 2026-03-30 10:00:00
type LogSamplingRule struct {
	Path     string        // 需要采样的路径
	Interval time.Duration // 采样间隔，在此时间内最多打印一次日志
}

// LogMiddlewareConfig 日志中间件配置
//
//	@author centonhuang
//	@update 2026-03-30 10:00:00
type LogMiddlewareConfig struct {
	SamplingRules []LogSamplingRule // 路径采样规则列表
}

// logSampler 日志采样器，记录每个路径的上次打印时间
type logSampler struct {
	mu       sync.Mutex
	lastLogs map[string]time.Time
}

func (s *logSampler) shouldLog(path string, interval time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if last, ok := s.lastLogs[path]; ok && now.Sub(last) < interval {
		return false
	}
	s.lastLogs[path] = now
	return true
}

// LogMiddleware 日志中间件
//
//	@param cfg LogMiddlewareConfig
//	@return fiber.Handler
//	@author centonhuang
//	@update 2026-03-30 10:00:00
func LogMiddleware(cfg LogMiddlewareConfig) fiber.Handler {
	samplingIndex := make(map[string]time.Duration, len(cfg.SamplingRules))
	for _, rule := range cfg.SamplingRules {
		samplingIndex[rule.Path] = rule.Interval
	}

	sampler := &logSampler{lastLogs: make(map[string]time.Time, len(cfg.SamplingRules))}

	return func(c *fiber.Ctx) error {
		start := time.Now().UTC()
		path := c.Path()
		query := string(c.Request().URI().QueryString())

		err := c.Next()

		// 对匹配采样规则的路径，按间隔控制日志频率（错误始终打印）
		if err == nil {
			if interval, ok := samplingIndex[path]; ok {
				if !sampler.shouldLog(path, interval) {
					return err
				}
			}
		}

		logger := logger.WithFCtx(c)

		latency := time.Since(start)

		fields := []zap.Field{
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
			zap.String("latency", latency.String()),
		}

		fields = append(fields, buildRequestHeadersFields(c)...)
		fields = append(fields, buildResponseHeadersFields(c)...)

		if strings.Contains(string(c.Request().Header.ContentType()), "application/json") {
			fields = append(fields, buildRequestBodyFields(c, logger)...)
		}

		// FIXME: get response body will break sse
		// reference: https://github.com/gofiber/fiber/issues/429
		// reference: https://github.com/samber/slog-fiber/issues/68
		if strings.Contains(string(c.Response().Header.ContentType()), "application/json") { // response header content-type is not text/event-stream
			fields = append(fields, buildResponseBodyFields(c, logger)...)
		}

		if err != nil {
			fields = append([]zap.Field{zap.Error(err)}, fields...)
			logger.Error("[LogMiddleware] error", fields...)
		} else {
			logger.Info("[LogMiddleware] info", fields...)
		}

		return err
	}
}

// buildRequestBodyFields 输出 JSON 请求体字段。
func buildRequestBodyFields(c *fiber.Ctx, logger *zap.Logger) []zap.Field {
	request := make(map[string]any)
	if reqBody := c.Body(); reqBody != nil {
		if err := sonic.Unmarshal(reqBody, &request); err != nil {
			logger.Warn("[LogMiddleware] unmarshal request error", zap.ByteString("request", reqBody), zap.Error(err))
		}
	}
	return []zap.Field{zap.Dict("request", lo.MapToSlice(request, func(key string, value any) zap.Field {
		return zap.Any(key, value)
	})...)}
}

// buildResponseBodyFields 输出 JSON 响应体字段（截断过长值）。
func buildResponseBodyFields(c *fiber.Ctx, logger *zap.Logger) []zap.Field {
	response := make(map[string]any)
	if respBody := c.Response().Body(); respBody != nil {
		if err := sonic.Unmarshal(respBody, &response); err != nil {
			logger.Warn("[LogMiddleware] unmarshal response error", zap.ByteString("response", respBody), zap.Error(err))
		}
	}
	truncated := util.TruncateMapValues(response, constant.LogFieldValueMaxLength)
	return []zap.Field{zap.Dict("response", lo.MapToSlice(truncated, func(key string, value any) zap.Field {
		return zap.Any(key, value)
	})...)}
}

// sensitiveHeadersForLog 需要掩码的敏感头列表。
var sensitiveHeadersForLog = []string{
	constant.HTTPHeaderAuthorization,
	constant.HTTPHeaderAPIKey,
	constant.HTTPHeaderCookie,
	constant.HTTPHeaderSetCookie,
}

func isSensitiveHeaderForLog(key string) bool {
	return lo.ContainsBy(sensitiveHeadersForLog, func(h string) bool { return strings.EqualFold(key, h) })
}

// buildRequestHeadersFields 输出请求头，敏感头（Authorization/API Key/Cookie）掩码为占位符。
func buildRequestHeadersFields(c *fiber.Ctx) []zap.Field {
	reqHeaders := make(map[string]any, len(c.GetReqHeaders()))
	for k, v := range c.GetReqHeaders() {
		if isSensitiveHeaderForLog(k) {
			reqHeaders[k] = constant.MaskSecretPlaceholder
			continue
		}
		reqHeaders[k] = v
	}
	return []zap.Field{zap.Dict("request-headers", lo.MapToSlice(reqHeaders, func(key string, value any) zap.Field {
		return zap.Any(key, value)
	})...)}
}

// buildResponseHeadersFields 输出响应头，Set-Cookie 等敏感头掩码为占位符。
func buildResponseHeadersFields(c *fiber.Ctx) []zap.Field {
	respHeaders := make(map[string]any, len(c.GetRespHeaders()))
	for k, v := range c.GetRespHeaders() {
		if isSensitiveHeaderForLog(k) {
			respHeaders[k] = constant.MaskSecretPlaceholder
			continue
		}
		respHeaders[k] = v
	}
	return []zap.Field{zap.Dict("response-headers", lo.MapToSlice(respHeaders, func(key string, value any) zap.Field {
		return zap.Any(key, value)
	})...)}
}
