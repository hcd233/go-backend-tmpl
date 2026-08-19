package constant

import "time"

const (
	// HTTPClientTimeout HTTP 客户端总超时时间。
	HTTPClientTimeout = 5 * time.Minute

	// HTTPDialTimeout HTTP 建连超时时间。
	HTTPDialTimeout = 10 * time.Second

	// HTTPKeepAlive HTTP keepalive 周期。
	HTTPKeepAlive = 30 * time.Second

	// HTTPTLSHandshakeTimeout TLS 握手超时时间。
	HTTPTLSHandshakeTimeout = 10 * time.Second

	// HTTPResponseHeaderTimeout 等待响应头超时时间。
	HTTPResponseHeaderTimeout = 30 * time.Second

	// HTTPIdleConnTimeout 空闲连接回收时间。
	HTTPIdleConnTimeout = 90 * time.Second
)

const (
	// HTTPMaxIdleConns 全局空闲连接上限。
	HTTPMaxIdleConns = 100

	// HTTPMaxIdleConnsPerHost 单 Host 空闲连接上限。
	HTTPMaxIdleConnsPerHost = 20
)

const (
	// ShutdownTimeout 整体优雅关闭最大超时（fx.StopTimeout 与 cmd 层共用）。
	ShutdownTimeout = 60 * time.Second

	// FiberShutdownTimeout HTTP server 关闭超时。
	FiberShutdownTimeout = 30 * time.Second

	// PoolStopTimeout 协程池停止超时。
	PoolStopTimeout = 30 * time.Second

	// CronStopTimeout 定时任务停止超时。
	CronStopTimeout = 30 * time.Second

	// InflightDrainSoftTimeout inflight 排空 soft 窗口：等待请求自然完成。
	InflightDrainSoftTimeout = 5 * time.Minute

	// InflightDrainHardTimeout inflight 排空 hard 窗口：广播取消后等待收尾。
	InflightDrainHardTimeout = 30 * time.Second
)

const (
	// HTTPHeaderAuthorization 认证头。
	HTTPHeaderAuthorization = "Authorization"

	// HTTPHeaderAPIKey API Key 头。
	HTTPHeaderAPIKey = "X-API-Key"

	// HTTPHeaderCookie Cookie 头。
	HTTPHeaderCookie = "Cookie"

	// HTTPHeaderSetCookie Set-Cookie 头。
	HTTPHeaderSetCookie = "Set-Cookie"
)
