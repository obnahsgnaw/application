package locker

import (
	"testing"
	"time"
	"wapi/config"
	"wapi/config/rds"
)

func TestNewRedisDistributedLocker(t *testing.T) {
	err := rds.InitRedis(&config.Redis{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "20210606123456",
	}, true)
	if err != nil {
		t.Error("redis init failed")
		return
	}
	l1, _ := NewRedisDistributedLocker(rds.Rds, "test1", 25*time.Second)

	if l1 == nil {
		t.Error("locker 1 need, but nil")
		return
	}

	l2, _ := NewRedisDistributedLocker(rds.Rds, "test1", 25*time.Second)
	if l2 != nil {
		t.Error("locker 2 nil need, but locker get", l2.val)
		return
	}

	l1.Unlock()

	l3, _ := NewRedisDistributedLocker(rds.Rds, "test1", 5*time.Second)
	if l3 == nil {
		t.Error("locker 3 need, but nil get")
		return
	}

	time.Sleep(5 * time.Second)

	l4, _ := NewRedisDistributedLocker(rds.Rds, "test1", 5*time.Second)
	if l4 == nil {
		t.Error("locker 4 need, but nil get")
		return
	}
}
