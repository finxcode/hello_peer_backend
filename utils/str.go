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
	s := strings.Trim(*strPtr, " ")
	if len(s) == 0 {
		return nil
	}
	return strings.Split(s, sep)
}

func ConcatImageUrl(filename, baseUrl string) string {
	return baseUrl + filename
}

func ConcatImagesUrl(filenames *[]string, baseUrl string) []string {
	if len(*filenames) == 0 {
		return nil
	}
	resStr := make([]string, len(*filenames))
	for i, filename := range *filenames {
		resStr[i] = ConcatImageUrl(filename, baseUrl)
	}

	return resStr
}
