package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

// Md5 生成32位MD5摘要
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// RandomString 生成长度为length的随机字符串
func RandomString(n int64) string {
	var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	letterLength := len(letters)
	for i := range result {
		result[i] = letters[rand.Intn(letterLength)]
	}
	return string(result)
}
