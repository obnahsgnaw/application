package locker

import (
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestNewRedisDistributedLocker(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "20210606123456",
		DB:       0,
	})
	l1, _ := NewRedisDistributedLocker(rds, "test1", 25*time.Second)

	if l1 == nil {
		t.Error("locker 1 need, but nil")
		return
	}

	l2, _ := NewRedisDistributedLocker(rds, "test1", 25*time.Second)
	if l2 != nil {
		t.Error("locker 2 nil need, but locker get", l2.val)
		return
	}

	l1.Unlock()

	l3, _ := NewRedisDistributedLocker(rds, "test1", 5*time.Second)
	if l3 == nil {
		t.Error("locker 3 need, but nil get")
		return
	}

	time.Sleep(5 * time.Second)

	l4, _ := NewRedisDistributedLocker(rds, "test1", 5*time.Second)
	if l4 == nil {
		t.Error("locker 4 need, but nil get")
		return
	}
}
