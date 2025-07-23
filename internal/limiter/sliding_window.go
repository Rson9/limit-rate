/*
滑动窗口限流算法
*/

package limiter

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu sync.Mutex
	// 存放每个时间窗口的请求数
	counts map[int64]int
	// 限制请求总数
	limit int64
	// 总时间窗口大小
	interval time.Duration
}

func NewSlidingWindowLimiter(limit int64, interval time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:    limit,
		interval: interval,
		counts:   make(map[int64]int),
	}
}

var _ Limter = (*SlidingWindowLimiter)(nil)

func (s *SlidingWindowLimiter) Check() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	windowStart := now - int64(s.interval.Seconds())

	var total int
	for ts, count := range s.counts {
		if ts < windowStart {
			delete(s.counts, ts)
		} else {
			total += count
		}
	}

	if int64(total) >= s.limit {
		return false
	}

	s.counts[now]++
	return true
}
