// Package pool 协程池管理器
//
//	author centonhuang
//	update 2026-02-04 16:10:57
package pool

import (
	"context"

	"github.com/alitto/pond/v2"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/dto"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
)

// Manager 协程池管理器，通过依赖注入构造。
//
//	author centonhuang
//	update 2026-01-31 16:00:00
type Manager struct {
	pingPool pond.Pool
}

// InitPoolManager 初始化协程池管理器并返回实例（供依赖注入使用）。
//
//	@return *Manager
//	@author centonhuang
//	@update 2026-01-31 03:37:28
func InitPoolManager() *Manager {
	return &Manager{
		pingPool: pond.NewPool(config.PoolWorkers, pond.WithQueueSize(config.PoolQueueSize)),
	}
}

// SubmitPingTask 提交异步 ping 示例任务。
//
//	@receiver pm *Manager
//	@param task
//	@return error
//	@author centonhuang
//	@update 2026-02-04 16:10:57
func (pm *Manager) SubmitPingTask(task *dto.PingTask) error {
	logger := logger.WithCtx(task.Ctx)
	return pm.pingPool.Go(func() {
		logger.Info("[PoolManager] async ping success")
	})
}

// StopWithContext 带超时优雅停止所有协程池（供优雅退出钩子调用）。
//
//	@receiver pm *Manager
//	@param ctx context.Context 停止超时上下文
//	@return error
func (pm *Manager) StopWithContext(ctx context.Context) error {
	if pm.pingPool == nil {
		return nil
	}
	task := pm.pingPool.Stop()
	done := make(chan struct{})
	go func() {
		_ = task.Wait() //nolint:errcheck // Stop future 无错误返回
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
