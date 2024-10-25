package security

import (
	"testing"
)

func TestAes(t *testing.T) {
	s := NewEsCrypto(Aes256, CbcMode)
	key := []byte("1234567890123456")
	data := `123dfg`
	b, iv, _ := s.Encrypt([]byte(data), key, true)
	b, _ = s.Decrypt(b, key, iv, true)
	println(string(b))
}
