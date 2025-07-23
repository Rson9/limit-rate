package limiter

import (
	"fmt"
	"time"
)

type Strategy string

const (
	FixedWindow   Strategy = "fixed"
	SlidingWindow Strategy = "sliding"
	LeakyBucket   Strategy = "leaky"
	TokenBucket   Strategy = "token"
	SlidingLog    Strategy = "log"
)

type Config struct {
	Strategy Strategy      `mapstructure:"strategy"`
	Limit    int64         `mapstructure:"limit"`
	Interval time.Duration `mapstructure:"interval"`
	Rate     int64         `mapstructure:"rate"`
}

func NewLimiter(cfg Config) (Limter, error) {
	switch cfg.Strategy {
	case FixedWindow:
		return NewFixed_windows_tryAcquire(cfg.Limit, cfg.Interval), nil
	case SlidingWindow:
		return NewSlidingWindowLimiter(cfg.Limit, cfg.Interval), nil
	case LeakyBucket:
		return NewLeakyBucketImproved(cfg.Limit, cfg.Rate), nil
	case TokenBucket:
		return NewTokenBucketImproved(cfg.Limit, cfg.Rate), nil
	case SlidingLog:
		return NewSlidingLog(int(cfg.Limit), cfg.Interval), nil
	default:
		return nil, fmt.Errorf("unsupported strategy: %s", cfg.Strategy)
	}
}
