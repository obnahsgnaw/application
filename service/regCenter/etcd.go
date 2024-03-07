package regCenter

import (
	"context"
	"github.com/obnahsgnaw/application/pkg/etcd"
	"github.com/obnahsgnaw/application/pkg/etcd/registercenter"
	"time"
)

type EtcdRegister struct {
	register *registercenter.EtcdRegister
}

func NewEtcdRegister(endpoints []string, opTimeout time.Duration) (*EtcdRegister, error) {
	r := &EtcdRegister{
		register: registercenter.New("", "", endpoints, opTimeout),
	}
	if err := r.register.Init(); err != nil {
		return nil, err
	}

	return r, nil
}

func (e *EtcdRegister) Release() {
	if e.register != nil {
		e.register.Release()
	}
}

func (e *EtcdRegister) Register(ctx context.Context, key, val string, ttl int64) error {
	if e.register == nil {
		return nil
	}
	return e.register.RegisterSimpleServer(ctx, key, val, ttl)
}

func (e *EtcdRegister) Unregister(ctx context.Context, key string) error {
	if e.register == nil {
		return nil
	}
	return e.register.Delete(ctx, key)
}

func (e *EtcdRegister) Watch(ctx context.Context, keyPrefix string, handler func(key string, val string, isDel bool)) error {
	if e.register != nil {
		e.register.WatchSimpleServer(ctx, keyPrefix, handler)
	}
	return nil
}

func (e *EtcdRegister) LastPrefixedIndex(ctx context.Context, keyPrefix string, indexParser func(key string) int) (int, error) {
	return etcd.GetLastIndex(ctx, e.Etcd().Conn(), keyPrefix, 5*time.Second, indexParser)
}

func (e *EtcdRegister) Etcd() *registercenter.EtcdRegister {
	return e.register
}
