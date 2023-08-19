package etcd

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func Puts(ctx context.Context, c *clientv3.Client, kvs map[string]string, leaseId clientv3.LeaseID) error {
	var opOptions []clientv3.OpOption
	var ops []clientv3.Op
	if leaseId > 0 {
		opOptions = append(opOptions, clientv3.WithLease(leaseId))
	}
	for k, v := range kvs {
		ops = append(ops, clientv3.OpPut(k, v, opOptions...))
	}
	ctx1, cl := context.WithTimeout(ctx, OpTtl)
	defer cl()
	txnRes, err := c.Txn(ctx1).Then(ops...).Commit()

	if err != nil {
		return err
	}
	if !txnRes.Succeeded {
		return errors.New("etcd error: tx put failed")
	}

	return nil
}

func PutsWithKeepalive(ctx context.Context, c *clientv3.Client, kvs map[string]string, leaseTtl int64, opTimeout time.Duration) error {
	lease, err := Grant(ctx, c, leaseTtl, opTimeout)
	if err != nil {
		return err
	}

	err = Puts(ctx, c, kvs, lease.ID)
	if err != nil {
		return err
	}

	return KeepAlive(ctx, c, lease.ID, func() {
		_ = PutsWithKeepalive(ctx, c, kvs, leaseTtl, opTimeout)
	})
}
