package etcd3

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Client struct {
	ctx context.Context
	c   *clientv3.Client
	ttl time.Duration
}

func New(ctx context.Context, endpoints []string, timeout time.Duration) (*Client, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		ctx: ctx,
		c:   c,
		ttl: time.Second * 5,
	}, nil
}

func (s *Client) Client() *clientv3.Client {
	return s.c
}

func (s *Client) Put(key, val string, leaseId clientv3.LeaseID) (*clientv3.PutResponse, error) {
	var options []clientv3.OpOption

	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()

	if leaseId > 0 {
		options = append(options, clientv3.WithLease(leaseId))
	}

	return s.c.Put(tx, key, val, options...)
}
func (s *Client) PutTtl(key, val string, ttl int64) (*clientv3.PutResponse, error) {
	var leaseId clientv3.LeaseID

	if ttl > 0 {
		lease, err := s.Grant(ttl)
		if err != nil {
			return nil, err
		}
		leaseId = lease.ID
	}

	return s.Put(key, val, leaseId)
}
func (s *Client) Get(key string) (*mvccpb.KeyValue, bool, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	resp, err := s.c.Get(tx, key)
	if err != nil {
		return nil, false, err
	}
	if resp.Count > 0 {
		return resp.Kvs[0], true, nil
	} else {
		return nil, false, nil
	}
}
func (s *Client) Delete(key string) (*mvccpb.KeyValue, bool, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()

	resp, err := s.c.Delete(tx, key)
	if err != nil {
		return nil, false, err
	}
	if resp.Deleted > 0 {
		return resp.PrevKvs[0], true, nil
	}
	return nil, false, nil
}
func (s *Client) Exist(key string) (bool, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	resp, err := s.c.Get(tx, key, clientv3.WithCountOnly())
	if err != nil {
		return false, err
	}
	return resp.Count > 0, nil
}
func (s *Client) Puts(kvs map[string]string, leaseId clientv3.LeaseID) (bool, error) {
	var options []clientv3.OpOption
	if leaseId > 0 {
		options = append(options, clientv3.WithLease(leaseId))
	}

	var ops []clientv3.Op
	for k, v := range kvs {
		ops = append(ops, clientv3.OpPut(k, v, options...))
	}

	resp, err := s.c.Txn(s.ctx).Then(ops...).Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}
func (s *Client) PutsTtl(kvs map[string]string, ttl int64) (bool, error) {
	var leaseId clientv3.LeaseID

	if ttl > 0 {
		lease, err := s.Grant(ttl)
		if err != nil {
			return false, err
		}
		leaseId = lease.ID
	}

	return s.Puts(kvs, leaseId)
}
func (s *Client) Gets(key string) ([]*mvccpb.KeyValue, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	resp, err := s.c.Get(tx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		return resp.Kvs, nil
	} else {
		return nil, nil
	}
}
func (s *Client) Deletes(key string) ([]*mvccpb.KeyValue, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()

	resp, err := s.c.Delete(tx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Deleted > 0 {
		return resp.PrevKvs, nil
	}
	return nil, nil
}
func (s *Client) Exists(key string) (bool, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	resp, err := s.c.Get(tx, key, clientv3.WithCountOnly(), clientv3.WithPrefix())
	if err != nil {
		return false, err
	}
	return resp.Count > 0, nil
}
func (s *Client) Count(key string) (int64, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	resp, err := s.c.Get(tx, key, clientv3.WithCountOnly(), clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Count, nil
}
func (s *Client) Grant(ttl int64) (*clientv3.LeaseGrantResponse, error) {
	tx, cl := context.WithTimeout(s.ctx, s.ttl)
	defer cl()
	return s.c.Grant(tx, ttl)
}
func (s *Client) Keepalive(leaseId clientv3.LeaseID, retry func() error) (err error) {
	var alive <-chan *clientv3.LeaseKeepAliveResponse

	if alive, err = s.c.KeepAlive(s.ctx, leaseId); err != nil {
		return
	}
	go func() {
		for {
			select {
			case aliveResp := <-alive:
				if aliveResp == nil {
					for {
						if err = retry(); err != nil {
							time.Sleep(time.Second * 2)
							continue
						}
						return
					}
				}
			case <-s.ctx.Done():
				return
			}
		}
	}()
	return
}
func (s *Client) GrantKeepalive(ttl int64, retry func() error) (err error) {
	var lease *clientv3.LeaseGrantResponse

	if lease, err = s.Grant(ttl); err != nil {
		return
	}

	return s.Keepalive(lease.ID, retry)
}
func (s *Client) PutKeepalive(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	var lease *clientv3.LeaseGrantResponse

	if lease, err = s.Grant(ttl); err != nil {
		return
	}

	if resp, err = s.Put(key, val, lease.ID); err != nil {
		return
	}

	err = s.Keepalive(lease.ID, func() error {
		_, err1 := s.PutKeepalive(key, val, ttl)
		return err1
	})

	return
}
func (s *Client) PutsKeepalive(kvs map[string]string, ttl int64) (resp bool, err error) {
	var lease *clientv3.LeaseGrantResponse

	if lease, err = s.Grant(ttl); err != nil {
		return
	}

	if resp, err = s.Puts(kvs, lease.ID); err != nil {
		return
	}

	err = s.Keepalive(lease.ID, func() error {
		_, err1 := s.PutsKeepalive(kvs, ttl)
		return err1
	})

	return
}
func (s *Client) Watch(key string, onPut func(v string), onDel func()) {
	// Fetch first
	kv, ok, _ := s.Get(key)
	if ok {
		onPut(string(kv.Value))
	}
	// then watch
	wch := s.c.Watch(s.ctx, key)
	go func() {
		for {
			select {
			case wrs := <-wch:
				for _, ev := range wrs.Events {
					if ev.Type == clientv3.EventTypePut {
						if onPut != nil {
							onPut(string(ev.Kv.Value))
						}
					} else if ev.Type == clientv3.EventTypeDelete {
						if onDel != nil {
							onDel()
						}
					}
				}

			case <-s.ctx.Done():
				return

			}
		}
	}()
}
func (s *Client) Watches(key string, onPut func(k, v string), onDel func(k string)) {
	// Fetch first
	kvs, _ := s.Gets(key)
	for _, kv := range kvs {
		if onPut != nil {
			onPut(string(kv.Key), string(kv.Value))
		}
	}

	// then watch
	wch := s.c.Watch(s.ctx, key, clientv3.WithPrefix())

	go func() {
		for {
			select {
			case wrs := <-wch:
				for _, ev := range wrs.Events {
					if ev.Type == clientv3.EventTypePut {
						if onPut != nil {
							onPut(string(ev.Kv.Key), string(ev.Kv.Value))
						}
					} else if ev.Type == clientv3.EventTypeDelete {
						if onDel != nil {
							onDel(string(ev.Kv.Key))
						}
					}
				}

			case <-s.ctx.Done():
				return

			}
		}
	}()
}
