package bootstrap

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/inflight"
	"github.com/hcd233/aris-api-tmpl/internal/cron"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/cache"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/database"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/metrics"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/pool"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type lifecycleParams struct {
	fx.In

	Lifecycle       fx.Lifecycle
	App             *fiber.App
	DB              *gorm.DB
	Cache           *redis.Client
	PoolManager     *pool.Manager
	InflightTracker *inflight.Tracker
	MetricsFlusher  *metrics.Flusher
	CronManager     *cron.CronManager
	ListenHost      string `name:"listenHost"`
	ListenPort      string `name:"listenPort"`
}

// registerLifecycleHooks 注册优雅退出钩子。
//
// fx 的 OnStop 钩子按注册顺序的逆序执行，因此此处按期望停止顺序的
// 逆序注册。期望停止顺序：
// Cron → Inflight → Pool → HTTP → Logger → DB → Redis
// （Inflight 先于 Pool：drain 期间协程池必须存活，被截断请求的落库才能完成）
func registerLifecycleHooks(params lifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			return cache.CloseCache()
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			return database.CloseDatabase()
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			return logger.Logger().Sync()
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listenAddr := fmt.Sprintf("%s:%s", params.ListenHost, params.ListenPort)
			go func() {
				if listenErr := params.App.Listen(listenAddr); listenErr != nil {
					logger.Logger().Error("[Server] HTTP server listen error", zap.Error(listenErr))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), constant.FiberShutdownTimeout)
			defer cancel()
			if err := params.App.ShutdownWithContext(shutdownCtx); err != nil {
				logger.Logger().Error("[Server] HTTP server shutdown error", zap.Error(err))
			}
			return nil
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			poolCtx, cancel := context.WithTimeout(context.Background(), constant.PoolStopTimeout)
			defer cancel()
			if err := params.PoolManager.StopWithContext(poolCtx); err != nil {
				logger.Logger().Warn("[Server] Pool stop error", zap.Error(err))
			}
			return nil
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.MetricsFlusher.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.MetricsFlusher.Stop()
			return nil
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			params.InflightTracker.Drain(constant.InflightDrainSoftTimeout, constant.InflightDrainHardTimeout)
			return nil
		},
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			cronCtx, cancel := context.WithTimeout(context.Background(), constant.CronStopTimeout)
			defer cancel()
			if err := params.CronManager.StopAll(cronCtx); err != nil {
				logger.Logger().Warn("[Server] Cron stop error", zap.Error(err))
			}
			return nil
		},
	})
}
