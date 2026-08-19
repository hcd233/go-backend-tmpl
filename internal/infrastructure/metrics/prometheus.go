// Package metrics 运行时指标采集基础设施：
// Prometheus registry + HTTP 指标采集器 + 周期性快照落 Redis（多实例聚合）。
package metrics

import (
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// NewRegistry 创建 Prometheus Registry，并注册 Go runtime / process 默认采集器。
//
// 默认采集器提供 go_goroutines / go_memstats_alloc_bytes / process_cpu_seconds_total，
// 是运行时大盘的 goroutine / heap / CPU 数据来源。
//
//	@return *prometheus.Registry
func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	return registry
}

var _ = constant.MetricServiceName
