package security

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

var (
	B64Encoding = base64.StdEncoding
	hexEncoding = HexEncoding{}
)

type Encoder interface {
	EncodeToString([]byte) string
	DecodeString(string) ([]byte, error)
}

type HexEncoding struct{}

func (HexEncoding) EncodeToString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func (HexEncoding) DecodeString(data string) ([]byte, error) {
	return hex.DecodeString(data)
}
