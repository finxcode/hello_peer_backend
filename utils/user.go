package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateHPId() string {
	captcha := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	return fmt.Sprintf("HP%s", captcha)
}
