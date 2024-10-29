package security

import (
	"testing"
)

func TestAes(t *testing.T) {
	s := NewEsCrypto(Aes256, CbcMode)
	key := []byte("1234567890123456")
	data := `
aldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfkaaldfhkashfaksdhfkjahdfka
`
	println("data len:", len(data))
	b, iv, _ := s.Encrypt([]byte(data), key, true)
	println("encrypt len:", len(b))
	println("encrypt data: ", string(b))
	b, _ = s.Decrypt(b, key, iv, true)
	println("decrypt len:", len(b))
	println("decrypt data:", string(b))
}
