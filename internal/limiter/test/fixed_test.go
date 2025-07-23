package limiter_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rson9/limit-rate/internal/limiter"
)

func TestFixedWindowLimiter(t *testing.T) {
	limiter := limiter.NewFixed_windows_tryAcquire(100, time.Second)

	// 100次请求应通过
	for i := 0; i < 100; i++ {
		if !limiter.Check() {
			t.Fatalf("expected pass at iteration %d", i)
		}
	}

	// 第101次应被限流
	if limiter.Check() {
		t.Fatalf("expected limit exceeded at 101")
	}

	// 等待窗口重置
	time.Sleep(time.Second)

	if !limiter.Check() {
		t.Fatalf("expected pass after window reset")
	}
}

func TestFixedWindowLimiter_Concurrent(t *testing.T) {
	limiter := limiter.NewFixed_windows_tryAcquire(100, time.Second)
	var passCount int64
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Check() {
				atomic.AddInt64(&passCount, 1)
			}
		}()
	}
	wg.Wait()

	if passCount != 100 {
		t.Fatalf("expected 100 passed, got %d", passCount)
	}
}
