package limiter

import (
	"sort"
	"sync"
	"time"
)

/*
滑动日志算法 (Sliding Log Algorithm)
这是一种精确但消耗内存的滑动窗口算法。
*/

type SlidingLogImporved struct {
	mu     sync.Mutex
	limit  int           // 限制请求数 (使用 int 更常见，除非你的 limit 超过 2^63)
	window time.Duration // 时间窗口大小
	log    []int64       // 使用 Unix 纳秒时间戳记录请求
}

var _ Limter = (*SlidingLogImporved)(nil)

// NewSlidingLog 创建一个滑动日志限流器
// limit: 时间窗口内的最大请求数
// window: 时间窗口的长度 (e.g., 1 * time.Minute)
func NewSlidingLog(limit int, window time.Duration) *SlidingLogImporved {
	return &SlidingLogImporved{
		limit:  limit,
		window: window,
		log:    make([]int64, 0, limit), // 预分配容量以提高性能
	}
}

// Check 限流检查器
func (s *SlidingLogImporved) Check() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 获取当前窗口的起始时间点
	now := time.Now().UnixNano()
	windowStart := now - s.window.Nanoseconds()

	// 2. 清理日志中所有已过期的时间戳
	// 我们需要找到第一个 > windowStart 的时间戳的索引
	// `sort.Search` 使用二分查找，比线性扫描更高效 (O(log N))
	firstValidIndex := sort.Search(len(s.log), func(i int) bool {
		return s.log[i] >= windowStart
	})

	// 通过 re-slice 高效地移除过期日志，无内存重分配
	s.log = s.log[firstValidIndex:]

	// 3. 检查当前窗口内的请求数量
	if len(s.log) < s.limit {
		// 未达到限制，记录当前请求并返回 true
		s.log = append(s.log, now)
		return true
	}

	// 已达到限制，返回 false
	return false
}
