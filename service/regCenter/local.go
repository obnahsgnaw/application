package regCenter

import (
	"context"
	"strings"
	"time"
)

type regVal struct {
	Value    string
	ttl      int64
	expireAt time.Time
}

type LocalRegister struct {
	ctx      context.Context
	data     map[string]regVal
	watchers map[string][]func(key string, val string, isDel bool)
}

func NewLocalRegister(ctx context.Context) (*LocalRegister, error) {
	r := &LocalRegister{
		ctx:      ctx,
		data:     make(map[string]regVal),
		watchers: make(map[string][]func(key string, val string, isDel bool)),
	}

	return r, nil
}

func (e *LocalRegister) Release() {
}

func (e *LocalRegister) Register(_ context.Context, key, val string, ttl int64) error {
	v := regVal{
		Value:    val,
		ttl:      ttl,
		expireAt: time.Time{},
	}
	if ttl > 0 {
		v.expireAt = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	e.data[key] = v
	e.notify(key, v.Value, false)
	return nil
}

func (e *LocalRegister) notify(key, val string, del bool) {
	for wk, wts := range e.watchers {
		if key == wk || strings.HasPrefix(key, wk) {
			for _, wt := range wts {
				wt(key, val, del)
			}
		}
	}
}
func (e *LocalRegister) Unregister(_ context.Context, key string) error {
	if v, ok := e.data[key]; ok {
		delete(e.data, key)
		e.notify(key, v.Value, true)
	}
	return nil
}

func (e *LocalRegister) Watch(_ context.Context, keyPrefix string, handler func(key string, val string, isDel bool)) error {
	e.watchers[keyPrefix] = append(e.watchers[keyPrefix], handler)
	return nil
}

func (e *LocalRegister) LastPrefixedIndex(_ context.Context, keyPrefix string, indexParser func(key string) int) (int, error) {
	index := -1
	for k := range e.data {
		if k == keyPrefix || strings.HasPrefix(k, keyPrefix) {
			i := indexParser(k)
			if i > index {
				index = i
			}
		}
	}
	return index, nil
}

func (e *LocalRegister) watch() {
	go func() {
		for {
			time.Sleep(time.Second * 3)
			select {
			case <-e.ctx.Done():
				return
			default:
				for k, v := range e.data {
					if v.ttl > 0 && v.expireAt.Before(time.Now()) {
						_ = e.Unregister(e.ctx, k)
					}
				}
			}
		}
	}()
}
