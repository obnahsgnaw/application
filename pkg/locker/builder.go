package locker

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type Builder interface {
	Compete(key string, ttl time.Duration) (Locker, error)
}

type RedisLockBuilder struct {
	client *redis.Client
}

func NewRedisLockBuilder(client *redis.Client) *RedisLockBuilder {
	return &RedisLockBuilder{client: client}
}

func (b *RedisLockBuilder) Compete(key string, timeout time.Duration) (Locker, error) {
	return NewRedisDistributedLocker(b.client, key, timeout)
}
