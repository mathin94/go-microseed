package cache

import (
	"context"
	"time"

	"microseed/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRedis(cfg *config.Config, log *zap.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		DialTimeout:  500 * time.Millisecond,
	})
	return rdb, nil
}

func RegisterHooks(lc fx.Lifecycle, rdb *redis.Client, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := rdb.Ping(ctx).Result()
			if err == nil {
				log.Info("redis connected")
			}
			return err
		},
		OnStop: func(ctx context.Context) error {
			log.Info("closing redis")
			return rdb.Close()
		},
	})
}
