package cron

import (
	"context"

	"github.com/google/uuid"
	"github.com/hcd233/go-backend-tmpl/internal/constant"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// ExampleCron 示例定时任务
//
//	@author centonhuang
//	@update 2025-09-30 16:09:05
type ExampleCron struct {
	cron *cron.Cron
}

// NewExampleCron 创建示例定时任务
//
//	@return Cron
//	@author centonhuang
//	@update 2025-09-30 16:09:09
func NewExampleCron() Cron {
	return &ExampleCron{
		cron: cron.New(
			cron.WithLogger(newCronLoggerAdapter("QuotaCron", logger.Logger())),
		),
	}
}

// Start 启动示例定时任务
//
//	@receiver c *ExampleCron
//	@return error
//	@author centonhuang
//	@update 2025-09-30 16:11:28
func (c *ExampleCron) Start() error {
	// debug set 10 seconds
	entryID, err := c.cron.AddFunc("every 10s", c.doSomething)
	// c.cron.AddFunc("daily", c.deliverQuotas)
	if err != nil {
		logger.Logger().Error("[ExampleCron] add func error", zap.Error(err))
		return err
	}

	logger.Logger().Info("[ExampleCron] add func success", zap.Int("entryID", int(entryID)))

	c.cron.Start()

	return nil
}

func (c *ExampleCron) doSomething() {
	ctx := context.WithValue(context.Background(), constant.CtxKeyTraceID, uuid.New().String())
	logger := logger.WithCtx(ctx)

	logger.Info("[ExampleCron] doSomething success")
}
