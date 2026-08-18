package constant

import "time"

const (
	// MetricServiceName 指标服务标识。
	MetricServiceName = "aris-api-tmpl"

	// MetricNamespaceHTTP HTTP 指标命名空间（最终指标名形如 http_request_duration_seconds）。
	MetricNamespaceHTTP = "http"

	// MetricNameRequestDuration 请求时延直方图（不含 namespace）。
	MetricNameRequestDuration = "request_duration_seconds"

	// MetricNameRequests HTTP 请求结果 counter（不含 namespace；完整名 http_requests_total）。
	MetricNameRequests = "requests_total"

	MetricRequestDurationHelp = "HTTP request latency in seconds"
	MetricRequestsHelp        = "HTTP business requests by result"

	// MetricLabelResult HTTP 请求结果 counter 的结果 label。
	MetricLabelResult = "result"

	// MetricFullRequestDuration flusher 从 registry.Gather() 抽取快照时用的完整指标名。
	MetricFullRequestDuration = "http_request_duration_seconds"
	MetricFullGoGoroutines    = "go_goroutines"
	MetricFullGoHeapAlloc     = "go_memstats_alloc_bytes"
	MetricFullProcessCPU      = "process_cpu_seconds_total"
	MetricFullHTTPRequests    = "http_requests_total"
)

// HTTPResult HTTP 请求结果 counter 的 result 枚举。
const (
	HTTPResultSuccess = "success"
	HTTPResultFailure = "failure"
)

// PrometheusRequestDurationBuckets 请求时延直方图桶（秒）。
var PrometheusRequestDurationBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75,
	1, 2.5, 5, 10, 15, 30, 60, 120, 300, 600, 1800,
}

const (
	// RuntimeMetricsFlushInterval 每个 pod 采集并写入快照的间隔（= 时序分辨率）。
	RuntimeMetricsFlushInterval = 5 * time.Second

	// RuntimeMetricsRetention Redis 中运行时快照的留存窗口。
	RuntimeMetricsRetention = 24 * time.Hour

	// RuntimeMetricsUnknownInstance hostname 获取失败时的兜底实例标识。
	RuntimeMetricsUnknownInstance = "unknown"

	// DecimalBase 十进制基数（strconv.FormatInt 用）。
	DecimalBase = 10
)
