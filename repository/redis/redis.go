package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-web-scaffold/settings"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return
	}

	return
}

func Close() {
	if err := rdb.Close(); err != nil {
		zap.L().Error("Redis close failed: ", zap.Error(err))
	}
}
