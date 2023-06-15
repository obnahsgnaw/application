package etcd

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c, err := NewClient([]string{"127.0.0.1:2379"}, time.Second*5)

	if err != nil {
		t.Error(err)
	}
	defer c.Close()
}
