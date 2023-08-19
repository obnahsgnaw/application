package locker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/obnahsgnaw/application/pkg/utils"
	"strconv"
	"time"
)

const dLockerKeyPrefix = "distributed-locker"

type Locker interface {
	Unlock()
}

type RedisDistributedLocker struct {
	key       string
	val       string
	createdAt time.Time
	ttl       time.Duration
	client    *redis.Client
}

func lockerErr(msg string) error {
	return utils.TitledError("redis locker error", msg, nil)
}

func NewRedisDistributedLocker(client *redis.Client, key string, ttl time.Duration) (*RedisDistributedLocker, error) {
	now := time.Now()
	val := strconv.FormatInt(now.UnixNano(), 10)
	key = dLockerKeyPrefix + ":" + key
	if rs := client.SetNX(context.Background(), key, val, ttl); rs.Err() != nil {
		return nil, lockerErr(rs.Err().Error())
	} else {
		if rs.Val() {
			return &RedisDistributedLocker{
				key:       key,
				val:       val,
				createdAt: now,
				ttl:       ttl,
				client:    client,
			}, nil
		}

		return nil, lockerErr("locker exists")
	}

}

func (l *RedisDistributedLocker) Unlock() {
	if l.createdAt.Add(l.ttl).Before(time.Now()) {
		return
	}
	sc := `if redis.call("get",KEYS[1]) == ARGV[1] then return redis.call("del",KEYS[1]) else return 0 end`
	_ = l.client.Eval(context.Background(), sc, []string{l.key}, l.val)
}
