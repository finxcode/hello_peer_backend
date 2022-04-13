package utils

import (
	"math/rand"
	"strings"
	"time"
)

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func ParseToArray(strPtr *string, sep string) []string {
	if len(*strPtr) == 0 {
		return nil
	}
	return strings.Split(*strPtr, sep)
}
