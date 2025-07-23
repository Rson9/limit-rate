package limiter

import (
	"sync"
	"time"
)

/*
漏桶算法
*/

type LeakyBucketImproved struct {
	mu           sync.Mutex
	limit        float64   // 桶的容量 (使用 float64)
	rate         float64   // 漏水速率(per second) (使用 float64)
	water        float64   // 当前水量 (使用 float64)
	lastLeakTime time.Time // 上一次漏水时间
}

var _ Limter = (*LeakyBucketImproved)(nil)

func NewLeakyBucketImproved(limit, rate int64) *LeakyBucketImproved {
	return &LeakyBucketImproved{
		limit:        float64(limit),
		rate:         float64(rate),
		water:        0,
		lastLeakTime: time.Now(),
	}
}

func (l *LeakyBucketImproved) Check() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// 1. 计算自上次漏水以来，应该漏掉多少水，这次使用 float64 进行精确计算
	// 无论上次和本次请求间隔多久，都重新计算漏水量
	elapsed := now.Sub(l.lastLeakTime)
	leakedWater := elapsed.Seconds() * l.rate

	// 2. 更新当前水量，并重置漏水起始时间
	// 这一步是关键：先减去漏掉的水，再更新时间戳
	if leakedWater > 0 {
		l.water = l.water - leakedWater
		if l.water < 0 {
			l.water = 0 // 水不能为负
		}
	}
	// 无论是否漏水，都将 lastLeakTime 更新为当前时间，为下次计算做准备
	l.lastLeakTime = now

	// 3. 检查桶是否还有空间容纳新请求（水量+1）
	if l.water+1.0 <= l.limit {
		// 有空间，水量加1
		l.water++
		return true
	}

	// 桶满了，拒绝请求
	return false
}
