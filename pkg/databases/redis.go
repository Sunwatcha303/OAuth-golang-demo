package databases

import (
	"context"
	"fmt"

	"github.com/Sunwatcha303/OAuth-golang-demo/configs"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg *configs.Configs) (*redis.Client, error) {
	ctx := context.Background()

	addr, err := utils.ConnectionUrlBuilder("redis", cfg)
	if err != nil {
		return nil, err
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   cfg.Redis.Database,
	})

	_, err = redisDb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return redisDb, nil
}
