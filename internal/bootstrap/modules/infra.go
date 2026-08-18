// Package modules 按职责拆分 fx 依赖模块。
package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/inflight"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/cache"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/database"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/httpclient"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/metrics"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/pool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// InfraModule 基础设施模块：数据库、缓存、协程池、inflight 跟踪器、运行时指标。
var InfraModule = fx.Module("infra",
	fx.Provide(
		NewDB,
		NewCache,
		NewPoolManager,
		NewInflightTracker,
		metrics.NewRegistry,
		NewHTTPCollector,
		NewMetricsMiddleware,
		NewRuntimeMetricsCache,
		NewMetricsFlusher,
	),
	fx.Invoke(InitHTTPClient),
)

// NewDB 初始化数据库连接。
func NewDB() *gorm.DB {
	return database.InitDatabase()
}

// NewCache 初始化 Redis 客户端。
func NewCache() *redis.Client {
	return cache.InitCache()
}

// NewPoolManager 初始化协程池管理器。
func NewPoolManager() *pool.Manager {
	return pool.InitPoolManager()
}

// NewInflightTracker 创建 inflight 排空跟踪器。
func NewInflightTracker() *inflight.Tracker {
	return inflight.NewTracker()
}

// InitHTTPClient 初始化通用 HTTP 客户端（fx.Invoke 惰性执行）。
func InitHTTPClient() {
	httpclient.InitHTTPClient()
}

// NewHTTPCollector 创建 HTTP 指标采集器。
func NewHTTPCollector(registry *prometheus.Registry) *metrics.HTTPCollector {
	return metrics.NewHTTPCollector(registry)
}

// NewMetricsMiddleware 创建指标采集中间件。
func NewMetricsMiddleware(collector *metrics.HTTPCollector) fiber.Handler {
	return collector.Middleware()
}

// NewRuntimeMetricsCache 创建运行时指标 Redis 存储。
func NewRuntimeMetricsCache(client *redis.Client) *cache.RuntimeMetricsCache {
	return cache.NewRuntimeMetricsCache(client)
}

// NewMetricsFlusher 创建运行时指标快照落库器。
func NewMetricsFlusher(registry *prometheus.Registry, store *cache.RuntimeMetricsCache) *metrics.Flusher {
	return metrics.NewFlusher(registry, store, constant.RuntimeMetricsFlushInterval, constant.RuntimeMetricsRetention)
}
