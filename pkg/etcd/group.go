package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type SingleService struct {
	ctx           context.Context
	c             *clientv3.Client
	kvs           map[string]string
	name          string
	ttl           int64
	lockTime      int
	opeTimeout    time.Duration
	checkDuration time.Duration
	leaseId       clientv3.LeaseID
}

func NewSingleService(ctx context.Context, c *clientv3.Client, name string, kvs map[string]string) *SingleService {
	return &SingleService{
		ctx:           ctx,
		c:             c,
		kvs:           kvs,
		name:          name,
		ttl:           5,
		lockTime:      10,
		opeTimeout:    5 * time.Second,
		checkDuration: 5 * time.Second,
	}
}

func (s *SingleService) SetTtl(ttl int64) {
	s.ttl = ttl
}
func (s *SingleService) SetLockTime(lockSecond int) {
	s.lockTime = lockSecond
}
func (s *SingleService) SetOpTimeout(timeout time.Duration) {
	s.opeTimeout = timeout
}
func (s *SingleService) SetCheckInterval(interval time.Duration) {
	s.checkDuration = interval
}

// RegisterSingletonService 将一组kv组成为一个单服务， 不存在即添加， 存在即监听 监听到不存在时 抢锁添加维护
func (s *SingleService) RegisterSingletonService() error {
	var err error
	err = s.checkSingleton(func(exist bool) {
		if exist {
			s.listenSingleton()
		} else {
			err = s.registerSingleton()
		}
	})

	return err
}

func (s *SingleService) checkSingleton(handle func(exist bool)) (err error) {
	// 检查是否存在
	var exists bool
	if exists, err = Exists(s.ctx, s.c, singletonStatusKey(s.name), s.opeTimeout); err != nil {
		return
	}
	handle(exists)
	return
}

func (s *SingleService) listenSingleton() {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				// 定时检查
				_ = s.checkSingleton(func(exist bool) {
					if !exist {
						// 不存在 抢锁
						if locker, release, err := GetLocker(s.c, lockerName(s.name), s.lockTime); err == nil {
							// 抢锁后 注册
							if err = locker.TryLock(s.ctx); err == nil {
								_ = s.registerSingleton()
								release()
							}
						}
					}
				})
				time.Sleep(s.checkDuration)
			}
		}
	}()
}

func (s *SingleService) registerSingleton() (err error) {
	var lease *clientv3.LeaseGrantResponse
	lease, err = GrantAndKeepalive(s.ctx, s.c, s.ttl, s.opeTimeout, func() {
		_ = s.registerSingleton()
	})
	if err != nil {
		return
	}

	if _, err = Put(s.ctx, s.c, singletonStatusKey(s.name), "1", lease.ID); err != nil {
		return
	}
	s.leaseId = lease.ID
	return Puts(s.ctx, s.c, s.kvs, lease.ID)
}

func (s *SingleService) IsHost(host string) bool {
	if v, err := Get(s.ctx, s.c, singletonStatusKey(s.name), s.opeTimeout); err == nil {
		if v.Count > 0 {
			return string(v.Kvs[0].Value) == host
		}
	}
	return false
}

func (s *SingleService) RefreshKv(k, v string) error {
	s.kvs[k] = v
	if s.leaseId > 0 {
		if _, err := Put(s.ctx, s.c, k, v, s.leaseId); err != nil {
			return err
		}
	}

	return nil
}

func singletonStatusKey(name string) string {
	return "/singleton-service/status/" + name
}

func lockerName(name string) string {
	return "/lockers/" + name
}
