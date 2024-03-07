package registercenter

import (
	"context"
	"github.com/obnahsgnaw/application/pkg/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type EtcdRegister struct {
	clusterId string
	prefix    string
	opTimeout time.Duration
	endpoints []string
	client    *clientv3.Client
}

func New(clusterId, prefix string, endpoints []string, opTimeout time.Duration) *EtcdRegister {
	return &EtcdRegister{
		clusterId: clusterId,
		prefix:    prefix,
		opTimeout: opTimeout,
		endpoints: endpoints,
		client:    nil,
	}
}

// Release etcd client
func (r *EtcdRegister) Release() {
	if r.client != nil {
		_ = r.client.Close()
	}
}

// Init etcd client
func (r *EtcdRegister) Init() (err error) {
	r.client, err = etcd.NewClient(r.endpoints, r.opTimeout)
	return
}

// Prefix return prefix
func (r *EtcdRegister) Prefix() string {
	return r.prefix
}

// ClusterId return clusterId
func (r *EtcdRegister) ClusterId() string {
	return r.clusterId
}

// OpeTimeout return operate timeout
func (r *EtcdRegister) OpeTimeout() time.Duration {
	return r.opTimeout
}

// Conn return register client
func (r *EtcdRegister) Conn() *clientv3.Client {
	return r.client
}

// EtcdKey etcd key
func (r *EtcdRegister) EtcdKey(key ...string) string {
	var keys []string
	if r.clusterId != "" {
		keys = append(keys, r.clusterId)
	}
	if r.prefix != "" {
		keys = append(keys, r.prefix)
	}
	if len(key) > 0 {
		keys = append(keys, key...)
	}
	if len(keys) == 0 {
		return ""
	}
	return strings.Join(keys, "/")
}

// Clear prefixed key
func (r *EtcdRegister) Clear(ctx context.Context, prefixedKey string) (deleted int64, err error) {
	resp, err := r.client.Delete(ctx, prefixedKey, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, nil
}

// Delete key
func (r *EtcdRegister) Delete(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, key)
	return err
}

// RegisterSimpleServer register simple server to etcd
func (r *EtcdRegister) RegisterSimpleServer(ctx context.Context, key string, val string, leaseTtl int64) error {
	return etcd.PutWithKeepalive(ctx, r.client, key, val, leaseTtl, r.opTimeout)
}

// WatchSimpleServer watch simpler server prefixed key
func (r *EtcdRegister) WatchSimpleServer(ctx context.Context, key string, handler func(key, val string, isDel bool)) {
	if handler != nil {
		etcd.Watch(ctx, r.client, key, true, func(e *clientv3.Event) {
			handler(string(e.Kv.Key), string(e.Kv.Value), false)
		}, func(e *clientv3.Event) {
			handler(string(e.Kv.Key), "", true)
		})
	}
}

// RegisterSingletonServer register a singleton server, only one server maintain the keys, other server watch and wait maintain
func (r *EtcdRegister) RegisterSingletonServer(ctx context.Context, name string, kvs map[string]string, leaseTtl int64) (err error) {
	id := r.clusterId + "-" + name
	s := etcd.NewSingleService(ctx, r.client, id, kvs)
	s.SetOpTimeout(r.opTimeout)
	s.SetTtl(leaseTtl)
	err = s.RegisterSingletonService()
	return
}
