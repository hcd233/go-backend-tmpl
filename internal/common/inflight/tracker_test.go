package inflight

import (
	"testing"
	"time"
)

func TestTracker_TrackUntrack(t *testing.T) {
	t.Parallel()

	tracker := NewTracker()
	if !tracker.Track() {
		t.Fatal("Track() returned false in running state")
	}
	tracker.Untrack()
	if tracker.IsDraining() {
		t.Fatal("IsDraining() returned true before Drain")
	}
}

func TestTracker_TrackRejectedAfterDrain(t *testing.T) {
	t.Parallel()

	tracker := NewTracker()
	// 排空开始后不再接受新请求
	tracker.state.Store(constantDraining)
	if tracker.Track() {
		t.Fatal("Track() returned true in draining state")
	}
}

func TestTracker_DrainWaitsForNaturalCompletion(t *testing.T) {
	t.Parallel()

	tracker := NewTracker()
	if !tracker.Track() {
		t.Fatal("Track() returned false")
	}

	done := make(chan bool, 1)
	go func() {
		done <- tracker.Drain(2*time.Second, 1*time.Second)
	}()

	// 请求在 soft 窗口内完成 → Drain 返回 true
	tracker.Untrack()

	select {
	case ok := <-done:
		if !ok {
			t.Fatal("Drain() returned false when all requests completed naturally")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Drain() did not return after untrack")
	}
}

func TestTracker_CancelOnDrain(t *testing.T) {
	t.Parallel()

	tracker := NewTracker()
	derived := tracker.CancelOnDrain(t.Context())

	// 排空广播后派生 ctx 被取消
	tracker.Drain(0, 1*time.Second)
	select {
	case <-derived.Done():
		// 期望取消
	case <-time.After(2 * time.Second):
		t.Fatal("CancelOnDrain ctx not canceled after drain broadcast")
	}
}

// constantDraining 直接置状态便于单测排空拒绝分支。
const constantDraining = int32(1)
