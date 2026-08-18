package modules

import (
	"github.com/hcd233/aris-api-tmpl/internal/cron"
	"go.uber.org/fx"
)

// CronModule 定时任务模块：管理器 + 任务注册。
var CronModule = fx.Module("cron",
	fx.Provide(NewCronManager),
	fx.Invoke(InitCronJobs),
)

// NewCronManager 创建定时任务管理器。
func NewCronManager() *cron.CronManager {
	return cron.NewCronManager()
}

// InitCronJobs 初始化并注册定时任务（fx.Invoke 惰性执行）。
func InitCronJobs(manager *cron.CronManager) {
	cron.InitCronJobs(manager)
}
