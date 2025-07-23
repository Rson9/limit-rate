package limiter_test

import (
	"sync"
	"testing"
	"time"

	"github.com/rson9/limit-rate/internal/limiter"
)

func TestSlidingWindowLimiter_Basic(t *testing.T) {
	limiter := limiter.NewSlidingWindowLimiter(5, 5*time.Second)

	// 前 5 次应通过
	for i := 0; i < 5; i++ {
		if !limiter.Check() {
			t.Errorf("expected allowed at iteration %d", i)
		}
	}

	// 第 6 次应被限流
	if limiter.Check() {
		t.Errorf("expected denied at 6th request")
	}

	// 等待窗口滑动，模拟时间前进（这部分可用于实际滑动观察）
	time.Sleep(6 * time.Second)

	// 应该可以重新通过
	if !limiter.Check() {
		t.Errorf("expected allowed after window slide")
	}
}

func TestSlidingWindowLimiter_Concurrent(t *testing.T) {
	limiter := limiter.NewSlidingWindowLimiter(50, 1*time.Second)

	var allowed int64
	var denied int64
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Check() {
				// 允许
				allowed++
			} else {
				// 拒绝
				denied++
			}
		}()
	}
	wg.Wait()

	if allowed > 50 {
		t.Errorf("allowed requests should not exceed 50, got %d", allowed)
	}
	if denied == 0 {
		t.Errorf("expected some denied requests")
	}
}
