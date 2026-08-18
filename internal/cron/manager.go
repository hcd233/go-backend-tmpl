package cron

import (
	"context"
	"sync"
)

// CronManager 定时任务注册表与生命周期管理器。
//
// 任务通过 Register 注册，优雅退出时由 StopAll(ctx) 带超时并发停止。
// 后续需要热重载（Redis pub/sub 广播）时在此扩展 StartListener。
type CronManager struct {
	mu    sync.Mutex
	crons []Cron
}

// NewCronManager 创建定时任务管理器。
func NewCronManager() *CronManager {
	return &CronManager{}
}

// Register 注册一个定时任务。
func (m *CronManager) Register(c Cron) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.crons = append(m.crons, c)
}

// StopAll 带超时并发停止全部定时任务。
func (m *CronManager) StopAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.crons) == 0 {
		return nil
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		var wg sync.WaitGroup
		for _, c := range m.crons {
			wg.Add(1)
			go func(cron Cron) {
				defer wg.Done()
				cron.Stop()
			}(c)
		}
		wg.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
