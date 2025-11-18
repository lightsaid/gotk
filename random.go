package gotk

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const minBytesPerString = 32

var (
	characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
	tkRand     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandomInt(min, max int64) int64 {
	return min + tkRand.Int63n(max-min+1)
}

func RandomString(num int) string {
	var sb strings.Builder
	var size = len(characters)
	for i := 0; i < num; i++ {
		s := characters[tkRand.Intn(size)]
		sb.WriteByte(s)
	}
	return sb.String()
}

func Bytes(n int) ([]byte, error) {
	bs := make([]byte, n)
	num, err := cryptoRand.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("Bytes: %w", err)
	}
	if num < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}

	return bs, nil
}

// SafeString 随机生产安全的字符串,n是字节数，少于minBytesPerString则设为minBytesPerString
func SafeString(n int) (string, error) {
	if n < minBytesPerString {
		n = minBytesPerString
	}

	bs, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("SafeString: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bs), nil
}
