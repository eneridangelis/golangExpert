package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	client *redis.Client
}

func NewRedisStrategy(host, port string) *RedisStrategy {
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})
	return &RedisStrategy{client: client}
}

func (rs *RedisStrategy) IncrementAndGet(key string, expirationSeconds int) (int64, error) {
	ctx := context.Background()
	count, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		err = rs.client.Expire(ctx, key, time.Duration(expirationSeconds)*time.Second).Err()
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (rs *RedisStrategy) Close() error {
	return rs.client.Close()
}
