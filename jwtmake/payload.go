package jwtmake

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lightsaid/gotk/random"
)

// JWTToken 负载数据
// 包含数据： uid 唯一标识，jwt.RegisteredClaims
type JWTPayload struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

// NewJWTPayload 创建一个Token Payload，
// 如果claims.ID不存在则默认生成一个
// 如果claims.ExpiresAt 不存在 默认15分钟
func NewJWTPayload(uid string, claims jwt.RegisteredClaims) *JWTPayload {
	if claims.ID == "" {
		claims.ID = strconv.FormatInt(time.Now().UnixNano(), 10) + random.RandomString(8)
	}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	return &JWTPayload{
		UID:              uid,
		RegisteredClaims: claims,
	}
}
