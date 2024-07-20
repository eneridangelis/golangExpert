package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})
}

func GetRateLimitKey(prefix, identifier string) string {
	return prefix + ":" + identifier
}

func IncrementRequestCount(client *redis.Client, key string, expiration int) (int64, error) {
	ctx := context.Background()
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		err = client.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
