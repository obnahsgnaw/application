package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/obnahsgnaw/application/pkg/security"
	"strconv"
	"time"
)

func GenLocalId(prefix string) string {
	return prefix + "_" + genMd5(strconv.FormatInt(time.Now().UnixNano(), 10)+security.RandAlphaNum(10))
}

func genMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
