package redisobs

import (
	"context"

	"github.com/go-redis/redis/v7"
)

// WrapRedisClient adds opentracing measurements for commands and returns cloned client
func WrapRedisClient(ctx context.Context, client *redis.Client) *redis.Client {
	rds := client.WithContext(ctx)
	hook := &redisHook{rds.Options()}
	rds.AddHook(hook)
	return rds
}
