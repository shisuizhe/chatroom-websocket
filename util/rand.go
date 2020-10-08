package util

import (
	"math/rand"
	"time"
)

var charRunes = []rune("0123456789abcdef")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandID 返回大小为n的随机字符串
func RandID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = charRunes[rand.Intn(len(charRunes))]
	}
	return string(b)
}
