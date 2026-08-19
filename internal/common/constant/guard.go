package constant

import "time"

const (
	// GuardStrikeThreshold 观察窗口内触发封禁的违规次数阈值。
	GuardStrikeThreshold = 5

	// GuardStrikeWindow 违规计数的观察窗口时长。
	GuardStrikeWindow = 1 * time.Minute

	// GuardBanDuration 触发封禁后的封禁时长。
	GuardBanDuration = 1 * time.Hour

	// ScannerStrikeKeyTemplate 路由扫描违规计数 key 模板（%s = IP）。
	ScannerStrikeKeyTemplate = "scanner:strike:%s"

	// ScannerBanKeyTemplate 路由扫描封禁 key 模板（%s = IP）。
	ScannerBanKeyTemplate = "scanner:ban:%s"
)
