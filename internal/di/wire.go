//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/rson9/limit-rate/internal/config"
	"github.com/rson9/limit-rate/internal/limiter"
)

func InitLimiter() (limiter.Limter, error) {
	panic(wire.Build(
		config.LoadLimiterConfig, // 返回 limiter.Config
		limiter.NewLimiter,
	))
}
