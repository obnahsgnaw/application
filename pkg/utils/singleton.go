package utils

import "sync"

var statics sync.Map

func GetSingleton(key string, insBuilder func() interface{}) interface{} {
	if v, ok := statics.Load(key); !ok {
		v = insBuilder()
		statics.Store(key, v)
		return v
	} else {
		return v
	}
}

type Singleton struct {
	sync.Once
	instance interface{}
	builder  func() interface{}
}

func NewSingleton(builder func() interface{}) *Singleton {
	return &Singleton{
		builder: builder,
	}
}

func (s *Singleton) Get() interface{} {
	s.Do(func() {
		s.instance = s.builder()
	})

	return s.instance
}
