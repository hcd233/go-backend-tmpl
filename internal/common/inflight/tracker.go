// Package inflight 提供并发请求跟踪与两阶段排空能力，
// 是 K8s 优雅退出（流量摘除）的基础设施。
package inflight

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"go.uber.org/zap"
)

// Tracker 跟踪进行中的请求，支持两阶段排空（Drain）。
type Tracker struct {
	wg         sync.WaitGroup
	state      atomic.Int32
	cancelCh   chan struct{}
	cancelOnce sync.Once
}

// NewTracker 创建排空跟踪器。
func NewTracker() *Tracker {
	t := &Tracker{}
	t.state.Store(constant.InflightStateRunning)
	t.cancelCh = make(chan struct{})
	return t
}

// Track 登记一个进行中的请求；排空开始后返回 false。
func (t *Tracker) Track() bool {
	if t.state.Load() == constant.InflightStateDraining {
		return false
	}
	t.wg.Add(1)
	if t.state.Load() == constant.InflightStateDraining {
		t.wg.Done()
		return false
	}
	return true
}

// Untrack 释放一个请求。
func (t *Tracker) Untrack() {
	t.wg.Done()
}

// Drain 两阶段排空：soft 窗口内等待所有请求自然完成；soft 到点广播取消信号
// （CancelOnDrain 派生的 ctx 被取消，使阻塞的上游读返回 context canceled），
// 再等 hard 窗口让被截断的请求写完错误帧、计量并 Untrack。
//
//	@param soft time.Duration 自然等待窗口
//	@param hard time.Duration 广播后的收尾窗口
//	@return bool 所有请求是否已释放（hard 超时返回 false，由 HTTP shutdown 兜底）
func (t *Tracker) Drain(soft, hard time.Duration) bool {
	t.state.Store(constant.InflightStateDraining)

	done := make(chan struct{})
	go func() {
		defer close(done)
		t.wg.Wait()
	}()

	select {
	case <-done:
		logger.Logger().Info("[Inflight] All inflight requests completed")
		return true
	case <-time.After(soft):
		t.broadcastCancel()
		logger.Logger().Warn("[Inflight] Drain soft deadline reached, canceling inflight requests",
			zap.Duration("softTimeout", soft))
		select {
		case <-done:
			logger.Logger().Info("[Inflight] All inflight requests completed after cancel")
			return true
		case <-time.After(hard):
			logger.Logger().Warn("[Inflight] Drain hard deadline reached, some requests may not have completed",
				zap.Duration("hardTimeout", hard))
			return false
		}
	}
}

// CancelOnDrain 返回 ctx 的派生 context：drain soft deadline 广播时取消派生 ctx，
// 使依赖该 ctx 的阻塞操作（如上游 SSE 读）退出。goroutine 在派生 ctx
// （随请求结束而 done）时退出，不泄漏。
func (t *Tracker) CancelOnDrain(ctx context.Context) context.Context {
	derived, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-t.cancelCh:
			cancel()
		case <-derived.Done():
		}
	}()
	return derived
}

func (t *Tracker) broadcastCancel() {
	t.cancelOnce.Do(func() { close(t.cancelCh) })
}

// IsDraining 是否处于排空状态。
func (t *Tracker) IsDraining() bool {
	return t.state.Load() == constant.InflightStateDraining
}
