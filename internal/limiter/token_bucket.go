package limiter

import (
	"sync"
	"time"
)

/*
令牌桶算法
*/

type TokenBucketImproved struct {
	// 容量
	limit int64
	// 放入令牌速率(per second)
	rate int64
	// 当前令牌数量
	currentTokens int64
	// 上次放入令牌时间
	lastTokenTime time.Time
	mu            sync.Mutex
}

// NewTokenBucketImproved 创建一个令牌桶
// limit: 令牌桶的容量
// rate: 放入令牌的速率
func NewTokenBucketImproved(limit, rate int64) *TokenBucketImproved {
	return &TokenBucketImproved{
		limit:         limit,
		rate:          rate,
		lastTokenTime: time.Now(),
		mu:            sync.Mutex{},
	}
}

var _ Limter = (*TokenBucketImproved)(nil)

func (t *TokenBucketImproved) Check() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 更新当前令牌数
	now := time.Now()
	elapsed := now.Sub(t.lastTokenTime)
	tokensToAdd := elapsed.Seconds() * float64(t.rate)
	t.currentTokens += int64(tokensToAdd)
	if t.currentTokens > t.limit {
		t.currentTokens = t.limit
	}
	t.lastTokenTime = now

	// 检查是否有足够的令牌
	if t.currentTokens > 0 {
		t.currentTokens--
		return true
	}

	return false
}
