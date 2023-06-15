package utils

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
)

func ProtoString(v proto.Message) (string, error) {
	if b, err := proto.Marshal(v); err == nil {
		return base64.StdEncoding.EncodeToString(b), nil
	} else {
		return "", err
	}
}

func ProtoOfString(b string, v proto.Message) error {
	bb, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return err
	}
	return proto.Unmarshal(bb, v)
}
