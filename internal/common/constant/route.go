package constant

const (
	// RoutePathRoot 根路径。
	RoutePathRoot = "/"

	// RoutePathHealth 健康检查路径（livenessProbe）。
	RoutePathHealth = "/health"

	// RoutePathReady 就绪检查路径（readinessProbe，draining 时返回 503）。
	RoutePathReady = "/ready"

	// RoutePathSSEHealth SSE 心跳检查路径。
	RoutePathSSEHealth = "/ssehealth"

	// RoutePathMetrics Prometheus 指标路径。
	RoutePathMetrics = "/metrics"
)

// InflightState 状态。
const (
	// InflightStateRunning 正常运行。
	InflightStateRunning = int32(0)

	// InflightStateDraining 排空中。
	InflightStateDraining = int32(1)
)
