package regCenter

import "context"

type None struct {
}

func NewNone() *None {
	return &None{}
}

func (s *None) Register(ctx context.Context, key, val string, ttl int64) error {
	return nil
}
func (s *None) Unregister(ctx context.Context, key string) error {
	return nil
}
func (s *None) Watch(ctx context.Context, keyPrefix string, handler func(key string, val string, isDel bool)) error {
	return nil
}
func (s *None) LastPrefixedIndex(ctx context.Context, keyPrefix string, indexParser func(key string) int) (int, error) {
	return 0, nil
}
