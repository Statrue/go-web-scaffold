package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			viper.GetString("redis.host"),
			viper.GetString("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.poolSize"),
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
