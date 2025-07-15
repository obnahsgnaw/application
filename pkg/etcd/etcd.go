package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

var etcdC *clientv3.Client
var OpTtl = time.Second * 5

// Client Get a client []string{"localhost:2379", "localhost:22379", "localhost:32379"}
func Client(endpoints []string, timeout time.Duration) (*clientv3.Client, error) {
	if etcdC == nil {
		c, err := NewClient(endpoints, timeout)
		if err != nil {
			return nil, err
		}

		etcdC = c
	}

	return etcdC, nil
}

// NewClient return a new client []string{"localhost:2379", "localhost:22379", "localhost:32379"}
func NewClient(endpoints []string, timeout time.Duration) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
}

func Put(ctx context.Context, c *clientv3.Client, key, val string, leaseId clientv3.LeaseID) (*clientv3.PutResponse, error) {
	var opOptions []clientv3.OpOption
	if leaseId > 0 {
		opOptions = append(opOptions, clientv3.WithLease(leaseId))
	}
	ctx1, cl := context.WithTimeout(ctx, OpTtl)
	defer cl()
	return c.Put(ctx1, key, val, opOptions...)
}

func Get(ctx context.Context, c *clientv3.Client, key string, opTimeout time.Duration, option ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx1, cl := context.WithTimeout(ctx, opTimeout)
	defer cl()
	return c.Get(ctx1, key, option...)
}

func Grant(ctx context.Context, c *clientv3.Client, ttl int64, opTtl time.Duration) (*clientv3.LeaseGrantResponse, error) {
	ctx1, cl := context.WithTimeout(ctx, opTtl)
	defer cl()
	return c.Grant(ctx1, ttl)
}

func KeepAlive(ctx context.Context, c *clientv3.Client, leaseId clientv3.LeaseID, _ func() error) (err error) {
	var alive <-chan *clientv3.LeaseKeepAliveResponse

	if alive, err = c.KeepAlive(ctx, leaseId); err != nil {
		return
	}
	go func() {
		for {
			select {
			case _, ok := <-alive:
				if !ok {
					return
				}
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 10):

			}
		}
	}()
	return
}

func GrantAndKeepalive(ctx context.Context, c *clientv3.Client, ttl int64, opTtl time.Duration, _ func() error) (*clientv3.LeaseGrantResponse, error) {
	lease, err := Grant(ctx, c, ttl, opTtl)
	if err != nil {
		return nil, err
	}
	err = KeepAlive(ctx, c, lease.ID, nil)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

func Watch(ctx context.Context, c *clientv3.Client, key string, prefixed bool, onPut func(e *clientv3.Event), onDel func(e *clientv3.Event)) {
	// Fetch first
	if prefixed {
		_ = GetPrefixed(ctx, c, key, OpTtl, func(kv *mvccpb.KeyValue) {
			if onPut != nil {
				onPut(&clientv3.Event{
					Type:   clientv3.EventTypePut,
					Kv:     kv,
					PrevKv: nil,
				})
			}
		})
	} else {
		resp, err := Get(ctx, c, key, OpTtl)
		if err != nil && resp != nil && resp.Count > 0 && onPut != nil {
			onPut(&clientv3.Event{
				Type:   clientv3.EventTypePut,
				Kv:     resp.Kvs[0],
				PrevKv: nil,
			})
		}
	}
	// then watch
	var opOptions []clientv3.OpOption
	if prefixed {
		opOptions = append(opOptions, clientv3.WithPrefix())
	}
	wch := c.Watch(ctx, key, opOptions...)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case wrs := <-wch:
				for _, ev := range wrs.Events {
					if ev.Type == clientv3.EventTypePut {
						if onPut != nil {
							onPut(ev)
						}
					} else if ev.Type == clientv3.EventTypeDelete {
						if onDel != nil {
							onDel(ev)
						}
					}
				}
			case <-time.After(time.Second * 10):

			}
		}
	}()
}

func PutWithKeepalive(ctx context.Context, c *clientv3.Client, key, val string, leaseTtl int64, opTimeout time.Duration) error {
	lease, err := Grant(ctx, c, leaseTtl, opTimeout)
	if err != nil {
		return err
	}

	_, err = Put(ctx, c, key, val, lease.ID)
	if err != nil {
		return err
	}

	return KeepAlive(ctx, c, lease.ID, nil)
}

func GetPrefixed(ctx context.Context, c *clientv3.Client, key string, opTimeout time.Duration, callback func(kv *mvccpb.KeyValue)) error {
	rs, err := Get(ctx, c, key, opTimeout, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if rs.Count > 0 {
		for _, kv := range rs.Kvs {
			callback(kv)
		}
	}

	return nil
}

func GetCount(ctx context.Context, c *clientv3.Client, prefix string, timeout time.Duration) (count int, err error) {
	var resp *clientv3.GetResponse
	if resp, err = Get(ctx, c, prefix, timeout, clientv3.WithPrefix(), clientv3.WithCountOnly()); err != nil {
		return
	}
	count = int(resp.Count)
	return
}

func GetLastIndex(ctx context.Context, c *clientv3.Client, prefix string, timeout time.Duration, indexParse func(key string) int) (index int, err error) {
	var resp *clientv3.GetResponse
	index = -1
	if resp, err = Get(ctx, c, prefix, timeout, clientv3.WithPrefix()); err != nil {
		return
	}
	if resp.Count == 0 {
		return
	}
	for _, kv := range resp.Kvs {
		i := indexParse(string(kv.Key))
		if i > index {
			index = i
		}
	}

	return
}

func Exists(ctx context.Context, c *clientv3.Client, key string, timeout time.Duration) (bool, error) {
	resp, err := Get(ctx, c, key, timeout, clientv3.WithCountOnly())
	if err != nil {
		return false, err
	}
	if resp.Count > 0 {
		return true, nil
	}
	return false, nil
}

func GetLocker(c *clientv3.Client, lockName string, ttl int) (locker *concurrency.Mutex, release func(), err error) {
	var session *concurrency.Session
	if session, err = concurrency.NewSession(c, concurrency.WithTTL(ttl)); err != nil {
		return
	}
	release = func() {
		_ = session.Close()
	}
	locker = concurrency.NewMutex(session, lockName)

	return
}
