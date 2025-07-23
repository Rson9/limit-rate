package limiter

import (
	"fmt"
	"sync"
	"time"
)

/*
固定窗口限流算法
现在是直接用的lock来防止并发
也可以对counter使用atomic来防止并发
*/

type Fixed_windows_tryAcquire struct {
	// 窗口大小
	count int64
	// 窗口时间
	time_windows time.Duration
	// 上次请求时间
	lastTime time.Time
	// 计数器
	counter int64
	// 防止并发
	lock sync.Mutex
}

var _ Limter = (*Fixed_windows_tryAcquire)(nil)

func NewFixed_windows_tryAcquire(count int64, time_windows time.Duration) *Fixed_windows_tryAcquire {
	return &Fixed_windows_tryAcquire{
		count:        count,
		time_windows: time_windows,
		lastTime:     time.Now(),
		counter:      0,
		lock:         sync.Mutex{},
	}
}

func (f *Fixed_windows_tryAcquire) Check() bool {
	fmt.Println("Fixed_windows_tryAcquire,count:", f.count)
	f.lock.Lock()
	defer f.lock.Unlock()
	// 获取系统当前时间
	now := time.Now()
	// 如果当前时间距离上次请求时间大于窗口时间
	if now.Sub(f.lastTime) > f.time_windows {
		// 重置计数器
		f.counter = 0
		// 更新请求时间
		f.lastTime = now
	}

	// 如果计数器小于窗口大小
	if f.counter < f.count {
		// 更新计数器
		f.counter++
		return true
	}
	return false
}
