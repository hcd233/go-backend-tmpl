package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// scannerGuardLua 路由扫描防护 Lua 脚本（原子操作）
//
// 当检测到一次路由未命中时调用：
//  1. 对 strike key 执行 INCR（违规计数 +1）
//  2. 若为首次记录，设置 strike key 的过期时间（观察窗口）
//  3. 若计数达到阈值，设置 ban key（封禁）并删除 strike key
//  4. 返回 [当前违规次数, 是否触发封禁(0/1)]
//
// KEYS[1]: strike key (scanner:strike:{ip})
// KEYS[2]: ban key    (scanner:ban:{ip})
// ARGV[1]: 封禁阈值
// ARGV[2]: 观察窗口 TTL（秒）
// ARGV[3]: 封禁时长 TTL（秒）
var scannerGuardLua = redis.NewScript(`
local strike_key = KEYS[1]
local ban_key = KEYS[2]
local threshold = tonumber(ARGV[1])
local window_ttl = tonumber(ARGV[2])
local ban_ttl = tonumber(ARGV[3])

local strikes = redis.call('INCR', strike_key)
if strikes == 1 then
    redis.call('EXPIRE', strike_key, window_ttl)
end

if strikes >= threshold then
    redis.call('SET', ban_key, '1', 'EX', ban_ttl)
    redis.call('DEL', strike_key)
    return {strikes, 1}
end

return {strikes, 0}
`)

// GuardConfig 路由扫描防护配置。
type GuardConfig struct {
	StrikeThreshold int           // 在观察窗口内触发封禁的违规次数阈值
	StrikeWindow    time.Duration // 违规计数的观察窗口时长
	BanDuration     time.Duration // 触发封禁后的封禁时长
	// IgnoredPaths 列出 404 时不计入违规的路径（如健康检查、浏览器常规探测）。
	// 为 nil 时不忽略任何路径。
	IgnoredPaths []string
	// AllowIPs 白名单 IP 列表，豁免路由扫描封禁。
	AllowIPs []string
}

// isRouteNotFound 判断 Fiber 返回的错误是否为路由未匹配。
func isRouteNotFound(err error) bool {
	var fiberErr *fiber.Error
	return errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound
}

// GuardMiddleware 路由扫描防护中间件。
//
// 在 Fiber 层拦截路由扫描行为：
//   - 请求到达时，检查 IP 是否已被封禁（Redis GET），若封禁则直接返回 403
//   - 请求处理后，若 Fiber 返回路由未命中错误（Cannot GET/POST/...），
//     则通过 Lua 脚本原子地记录违规并在达到阈值时自动封禁
func GuardMiddleware(cache *redis.Client, cfg GuardConfig) fiber.Handler {
	thresholdStr := strconv.Itoa(cfg.StrikeThreshold)
	windowTTLStr := strconv.FormatInt(int64(cfg.StrikeWindow.Seconds()), 10)
	banTTLStr := strconv.FormatInt(int64(cfg.BanDuration.Seconds()), 10)
	ignoredPaths := lo.SliceToMap(cfg.IgnoredPaths, func(p string) (string, struct{}) {
		return p, struct{}{}
	})
	allowIPs := lo.SliceToMap(cfg.AllowIPs, func(p string) (string, struct{}) {
		return p, struct{}{}
	})
	return func(c *fiber.Ctx) error {
		if cache == nil {
			logger.WithFCtx(c).Warn("[GuardMiddleware] Redis dependency is nil")
			return c.Next()
		}
		ip := c.IP()
		if _, allowed := allowIPs[ip]; allowed {
			return c.Next()
		}
		banKey := fmt.Sprintf(constant.ScannerBanKeyTemplate, ip)
		ctx := c.Context()

		banned, err := cache.Exists(ctx, banKey).Result()
		if err != nil {
			logger.WithFCtx(c).Warn("[GuardMiddleware] Failed to check ban status", zap.String("ip", ip), zap.Error(err))
		}
		if banned > 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		nextErr := c.Next()

		if isRouteNotFound(nextErr) {
			// 常见浏览器/爬虫对健康路径的探测不计入路由扫描违规。
			if _, skip := ignoredPaths[c.Path()]; skip {
				return nextErr
			}
			strikeKey := fmt.Sprintf(constant.ScannerStrikeKeyTemplate, ip)

			result, luaErr := scannerGuardLua.Run(
				ctx, cache,
				[]string{strikeKey, banKey},
				thresholdStr, windowTTLStr, banTTLStr,
			).Int64Slice()
			if luaErr != nil {
				logger.WithFCtx(c).Warn("[GuardMiddleware] Failed to execute strike script", zap.String("ip", ip), zap.Error(luaErr))
				return nextErr
			}

			if result[1] == 1 {
				logger.WithFCtx(c).Warn("[GuardMiddleware] IP banned due to route scanning",
					zap.String("ip", ip),
					zap.Int64("strikes", result[0]),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()))
			}
		}

		return nextErr
	}
}
