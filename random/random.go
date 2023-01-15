package random

import (
	"math/rand"
	"strings"
	"time"
)

const stringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	x := len(stringSource)

	for i := 0; i < n; i++ {
		s := stringSource[rand.Intn(x)]
		sb.WriteByte(s)
	}

	return sb.String()
}
