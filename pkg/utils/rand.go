package utils

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

// RandAlphaNum  生成随机字符串([a~zA~Z0~9])
func RandAlphaNum(lenNum int) string {
	str := strings.Builder{}
	length := len(chars)
	for i := 0; i < lenNum; i++ {
		l := chars[rand.Intn(length)]
		str.WriteString(l)
	}
	return str.String()
}

// RandNum  生成随机数字字符串([0~9])
func RandNum(lenNum int) string {
	str := strings.Builder{}
	length := 10
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[52+rand.Intn(length)])
	}
	return str.String()
}

// RandAlpha  生成随机字符串(a~zA~Z])
func RandAlpha(lenNum int) string {
	str := strings.Builder{}
	length := 52
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[rand.Intn(length)])
	}
	return str.String()
}

func RandInt(l int) int {
	return rand.Intn(l)
}
