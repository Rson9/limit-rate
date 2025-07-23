package config

import (
	"time"

	"github.com/rson9/limit-rate/internal/limiter"
	"github.com/spf13/viper"
)

func LoadLimiterConfig() (limiter.Config, error) {
	viper.SetConfigFile("configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return limiter.Config{}, err
	}

	var cfg limiter.Config
	if err := viper.UnmarshalKey("limiter", &cfg); err != nil {
		return limiter.Config{}, err
	}

	if cfg.Interval == 0 {
		cfg.Interval = time.Second
	}
	return cfg, nil
}
