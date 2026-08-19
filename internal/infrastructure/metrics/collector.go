package metrics

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// HTTPCollector 自实现的 HTTP 运行时指标采集器（替代 fiberprometheus）。
//
// 维护两项进程内指标：请求时延 histogram（请求总量 QPS 由 histogram 的 sample count 派生），
// 以及请求结果 counter（success/failure，Success Rate 数据源）。
type HTTPCollector struct {
	duration prometheus.Histogram
	requests *prometheus.CounterVec
	skipURIs map[string]struct{}
}

// NewHTTPCollector 在 registry 上注册 HTTP 指标并返回采集器。
func NewHTTPCollector(registry *prometheus.Registry) *HTTPCollector {
	duration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: constant.MetricNamespaceHTTP,
		Name:      constant.MetricNameRequestDuration,
		Help:      constant.MetricRequestDurationHelp,
		Buckets:   constant.PrometheusRequestDurationBuckets,
	})
	registry.MustRegister(duration)

	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: constant.MetricNamespaceHTTP,
			Name:      constant.MetricNameRequests,
			Help:      constant.MetricRequestsHelp,
		},
		[]string{constant.MetricLabelResult},
	)
	registry.MustRegister(requests)
	// 预置 success/failure 子序列：CounterVec 在首次 WithLabelValues 之前
	// 不产生任何序列，否则无流量时 Gather 不输出对应指标，快照恒为空。
	for _, result := range []string{constant.HTTPResultSuccess, constant.HTTPResultFailure} {
		requests.WithLabelValues(result).Add(0)
	}

	skipURIs := map[string]struct{}{
		constant.RoutePathHealth:    {},
		constant.RoutePathReady:     {},
		constant.RoutePathSSEHealth: {},
		constant.RoutePathMetrics:   {},
	}
	return &HTTPCollector{duration: duration, requests: requests, skipURIs: skipURIs}
}

// Middleware 返回记录请求时延与成功/失败计数的 Fiber 中间件。
//
// 须全局挂载（app.Use(mw)）以覆盖所有业务路由；探活与指标路径被跳过。
func (hc *HTTPCollector) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, skip := hc.skipURIs[c.Path()]; skip {
			return c.Next()
		}
		start := time.Now()
		err := c.Next()
		hc.duration.Observe(time.Since(start).Seconds())
		if c.Response().StatusCode() == fiber.StatusOK {
			hc.requests.WithLabelValues(constant.HTTPResultSuccess).Inc()
		} else {
			hc.requests.WithLabelValues(constant.HTTPResultFailure).Inc()
		}
		return err
	}
}
