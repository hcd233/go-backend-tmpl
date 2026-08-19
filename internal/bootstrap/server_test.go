package bootstrap

import (
	"testing"

	"github.com/hcd233/aris-api-tmpl/internal/common/inflight"
	"github.com/hcd233/aris-api-tmpl/internal/cron"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/pool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// buildTestApp 构造不依赖真实外部服务（DB/Redis）的 fx 应用，
// 用于验证 DI 图完整性、中间件与路由装配。
func buildTestApp(t *testing.T) *fx.App {
	t.Helper()

	customizers := []fx.Option{
		fx.NopLogger,
		fx.Replace(&gorm.DB{}),
		fx.Replace(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})),
		fx.Replace(pool.InitPoolManager()),
		fx.Replace(cron.NewCronManager()),
		fx.Replace(inflight.NewTracker()),
	}

	app := BuildFxApp("localhost", "0", customizers...)
	if err := app.Err(); err != nil {
		t.Fatalf("BuildFxApp() error = %v", err)
	}
	return app
}

func TestBuildFxApp_GraphIsConstructible(t *testing.T) {
	t.Parallel()

	buildTestApp(t)
}

func TestBuildFxApp_CreatesIsolatedApps(t *testing.T) {
	t.Parallel()

	first := buildTestApp(t)
	second := buildTestApp(t)
	if first == second {
		t.Fatal("BuildFxApp() reused fx.App instance")
	}
}
