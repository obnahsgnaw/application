package regCenter

import (
	"context"
	"github.com/obnahsgnaw/application/regtype"
	"strings"
)

type Register interface {
	Register(ctx context.Context, key, val string, ttl int64) error
	Unregister(ctx context.Context, key string) error
	Watch(ctx context.Context, keyPrefix string, handler func(key string, val string, isDel bool)) error
	LastPrefixedIndex(ctx context.Context, keyPrefix string, indexParser func(key string) int) (int, error)
}

type ServerInfo struct {
	Id      string
	Name    string
	Type    string
	EndType string
}

// RegInfo register info to register center
type RegInfo struct {
	AppId      string
	RegType    regtype.RegType
	ServerInfo ServerInfo
	Host       string
	Val        string // 单个值设置这个
	Ttl        int64
	KeyPreGen  RegKeyPrefixGenerator
	Values     map[string]string // 多个值设置这个
}

func (r *RegInfo) Prefix() string {
	if r.KeyPreGen == nil {
		r.KeyPreGen = DefaultRegKeyPrefixGenerator()
	}
	return r.KeyPreGen(r)
}
func (r *RegInfo) Key() string {
	prefix := r.Prefix()
	return strings.TrimPrefix(strings.Join([]string{prefix, r.ServerInfo.Id, r.Host}, "/"), "/")
}

func (r *RegInfo) Kvs() map[string]string {
	kvs := make(map[string]string)
	if r.Values != nil && len(r.Values) > 0 {
		key := r.Key()
		for k, v := range r.Values {
			kvs[key+"/"+k] = v
		}
	} else {
		kvs[r.Key()] = r.Val
	}
	return kvs
}

// RegKeyPrefixGenerator register key prefix generator  :  prefix/server-id/host
type RegKeyPrefixGenerator func(info *RegInfo) string

// 集群id/注册类型/前后台类型/服务类型/模块/host => host

// http注册
// dev/http/backend/api/auth/127.0.0.1:80 => 127.0.0.1:80

// rpc注册
// dev/rpc/backend/api/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/rpc/backend/tcp/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/rpc/backend/wss/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/rpc/backend/udp/auth/127.0.0.1:80 => 127.0.0.1:80

// doc注册
// dev/doc/backend/api/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/doc/backend/tcp/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/doc/backend/wss/auth/127.0.0.1:80 => 127.0.0.1:80
// dev/doc/backend/udp/auth/127.0.0.1:80 => 127.0.0.1:80

// DefaultRegKeyPrefixGenerator the default register generator:app-id/end-type/server-type/server-id/host => addr
func DefaultRegKeyPrefixGenerator() RegKeyPrefixGenerator {
	return func(info *RegInfo) string {
		return strings.Join([]string{info.AppId, info.RegType.String(), info.ServerInfo.EndType, info.ServerInfo.Type}, "/")
	}
}
func ActionRegKeyPrefixGenerator() RegKeyPrefixGenerator {
	return func(info *RegInfo) string {
		return strings.Join([]string{info.AppId, info.RegType.String(), info.ServerInfo.EndType, info.ServerInfo.Type, "action"}, "/")
	}
}
