// Package cron 定时任务模块
//
//	update 2024-12-09 15:55:25
package cron

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// Cron 定时任务接口
//
//	@author centonhuang
//	@update 2026-03-23 10:00:00
type Cron interface {
	Start() error
	Stop()
}

// InitCronJobs 初始化定时任务并注册到管理器。
//
//	@param manager *CronManager 任务管理器（优雅退出时按 StopAll 停止）
//	@author centonhuang
//	@update 2026-03-23 10:00:00
func InitCronJobs(manager *CronManager) {
	exampleCron := NewExampleCron()
	lo.Must0(exampleCron.Start())
	manager.Register(exampleCron)

	logger.Logger().Info("[Cron] Init cron jobs")
}

type cronLoggerAdapter struct {
	module string
	logger *zap.Logger
}

func newCronLoggerAdapter(module string, logger *zap.Logger) cronLoggerAdapter {
	if module == "" {
		module = "Cron"
	}
	module = strings.TrimSpace(strings.TrimRight(strings.TrimLeft(strings.TrimSpace(module), "["), "]"))
	return cronLoggerAdapter{module: module, logger: logger}
}

func (l cronLoggerAdapter) Error(err error, msg string, keysAndValues ...any) {
	zapKeyValues := []zap.Field{zap.Error(err)}
	zapKeyValues = append(zapKeyValues, convertZapKeyValues(keysAndValues...)...)
	l.logger.Error(fmt.Sprintf("[%s] %s", l.module, capitalizeFirst(msg)), zapKeyValues...)
}

func (l cronLoggerAdapter) Info(msg string, keysAndValues ...any) {
	zapKeyValues := convertZapKeyValues(keysAndValues...)
	l.logger.Info(fmt.Sprintf("[%s] %s", l.module, capitalizeFirst(msg)), zapKeyValues...)
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	for i, r := range s {
		return string(unicode.ToUpper(r)) + s[i+utf8.RuneLen(r):]
	}
	return s
}

func convertZapKeyValues(keysAndValues ...any) []zap.Field {
	if len(keysAndValues)%2 != 0 {
		panic("keysAndValues must be a slice of key-value pairs")
	}
	kvLen := len(keysAndValues) / 2
	zapKeyValues := make([]zap.Field, 0, kvLen)
	for i := range kvLen {
		key, value := keysAndValues[i*2].(string), keysAndValues[i*2+1]
		zapKeyValues = append(zapKeyValues, zap.Any(key, value))
	}
	return zapKeyValues
}
