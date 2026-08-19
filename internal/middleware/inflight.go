package middleware

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/common/inflight"
	"github.com/hcd233/aris-api-tmpl/internal/dto"
	"github.com/hcd233/aris-api-tmpl/internal/i18n"
)

var inflightHealthCheckPaths = map[string]struct{}{
	constant.RoutePathHealth:    {},
	constant.RoutePathReady:     {},
	constant.RoutePathSSEHealth: {},
}

// InflightMiddleware 并发请求跟踪中间件：登记进行中的请求，
// 优雅退出排空期间拒绝新请求并返回 503（K8s/负载均衡据此停止转发）。
// 503 是流量摘除信号，不属于业务错误语义，不受统一 200 契约约束。
func InflightMiddleware(tracker *inflight.Tracker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, skip := inflightHealthCheckPaths[c.Path()]; skip {
			return c.Next()
		}

		if !tracker.Track() {
			body, _ := sonic.Marshal(&dto.CommonRsp{ //nolint:errcheck // 静态结构体序列化不会失败
				Error: ierr.ErrInternal.BizError().Localize(i18n.FromCtx(c.UserContext())),
			})
			return c.Status(fiber.StatusServiceUnavailable).Type("json").Send(body)
		}
		defer tracker.Untrack()
		return c.Next()
	}
}
