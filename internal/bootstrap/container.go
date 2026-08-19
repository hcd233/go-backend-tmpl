// Package bootstrap wires application startup dependencies.
package bootstrap

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/hcd233/aris-api-tmpl/internal/api"
	"github.com/hcd233/aris-api-tmpl/internal/bootstrap/modules"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/inflight"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/enum"
	"github.com/hcd233/aris-api-tmpl/internal/handler"
	"github.com/hcd233/aris-api-tmpl/internal/middleware"
	"github.com/hcd233/aris-api-tmpl/internal/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// BuildFxAppOptions 组装应用依赖选项；customizers 供测试等场景追加覆盖。
func BuildFxAppOptions(host, port string, customizers ...fx.Option) []fx.Option {
	opts := []fx.Option{
		fx.Supply(
			fx.Annotate(host, fx.ResultTags(`name:"listenHost"`)),
			fx.Annotate(port, fx.ResultTags(`name:"listenPort"`)),
		),
		modules.InfraModule,
		modules.CronModule,
		modules.RepositoryModule,
		modules.ApplicationModule,
		modules.HandlerModule,
		fx.Provide(
			api.NewFiberApp,
			api.NewHumaAPI,
		),
		fx.Invoke(
			registerMiddlewares,
			registerRoutes,
			registerLifecycleHooks,
		),
		fx.StopTimeout(constant.ShutdownTimeout),
	}
	opts = append(opts, customizers...)
	return opts
}

// BuildFxApp 构建 fx 应用。
func BuildFxApp(host, port string, customizers ...fx.Option) *fx.App {
	return fx.New(BuildFxAppOptions(host, port, customizers...)...)
}

type middlewareParams struct {
	fx.In

	App               *fiber.App
	Cache             *redis.Client
	InflightTracker   *inflight.Tracker
	Registry          *prometheus.Registry
	MetricsMiddleware fiber.Handler
}

// registerMiddlewares 注册全局中间件链，顺序：Recover → Metrics → Inflight → Guard → Fgprof → CORS → Compress → Trace → Locale → Log
func registerMiddlewares(params middlewareParams) {
	// 标准 Prometheus 文本端点，供 Prometheus 抓取
	params.App.Get(constant.RoutePathMetrics, adaptor.HTTPHandler(promhttp.HandlerFor(params.Registry, promhttp.HandlerOpts{})))

	params.App.Use(
		middleware.RecoverMiddleware(),
		params.MetricsMiddleware,
		middleware.InflightMiddleware(params.InflightTracker),
		middleware.GuardMiddleware(params.Cache, middleware.GuardConfig{
			StrikeThreshold: constant.GuardStrikeThreshold,
			StrikeWindow:    constant.GuardStrikeWindow,
			BanDuration:     constant.GuardBanDuration,
			AllowIPs:        config.GuardAllowIPs,
			IgnoredPaths: []string{
				constant.RoutePathRoot,
				constant.RoutePathHealth,
				constant.RoutePathReady,
				constant.RoutePathSSEHealth,
				constant.RoutePathMetrics,
			},
		}),
		middleware.FgprofMiddleware(),
		middleware.CORSMiddleware(),
		middleware.CompressMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LocaleMiddleware(),
		middleware.LogMiddleware(middleware.LogMiddlewareConfig{
			SamplingRules: []middleware.LogSamplingRule{
				{Path: constant.RoutePathHealth, Interval: 5 * time.Minute},
				{Path: constant.RoutePathReady, Interval: 5 * time.Minute},
				{Path: constant.RoutePathSSEHealth, Interval: 5 * time.Minute},
			},
		}),
	)
}

type routeParams struct {
	fx.In

	App             *fiber.App
	HumaAPI         huma.API
	InflightTracker *inflight.Tracker
	PingHandler     handler.PingHandler
	TokenHandler    handler.TokenHandler
	Oauth2Handler   handler.Oauth2Handler
	UserHandler     handler.UserHandler
}

// registerRoutes 注册文档和 API 路由。
func registerRoutes(params routeParams) {
	// readiness 探针：draining（优雅退出排空）期间返回 503，让 K8s 先摘流量再停进程
	params.App.Get(constant.RoutePathReady, func(c *fiber.Ctx) error {
		if params.InflightTracker.IsDraining() {
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}
		return c.SendString("ok")
	})

	if config.Env != enum.EnvProduction {
		router.RegisterDocsRouter(params.App)
	}
	router.RegisterAPIRouter(params.HumaAPI, router.APIRouterDependencies{
		PingHandler:   params.PingHandler,
		TokenHandler:  params.TokenHandler,
		Oauth2Handler: params.Oauth2Handler,
		UserHandler:   params.UserHandler,
	})
}
