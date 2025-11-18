package gotk

import (
	"errors"
	"fmt"
	"time"

	"github.com/pascaldekloe/jwt"

	"github.com/google/uuid"
)

const minSecretKeySize = 32

var (
	ErrIssuerRequired = errors.New("issuer required")
	ErrSecretKeySize  = fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	ErrInvalidToken   = errors.New("token is invalid")
	ErrExpiredToken   = errors.New("token has expired")
)

type TokenPayload struct {
	// Token 唯一标识
	ID string `json:"id"`

	// 用户存储的数据
	Data string `json:"Data"`

	// 签发时间
	IssuedAt time.Time `json:"issuedAt"`

	// 过期时间
	ExpiredAt time.Time `json:"expiredAt"`
}

// NewTokenPayload 创建一个 TokenPayload 对象，
// 因为使用了 uuid.NewString() 可能会panic
func NewTokenPayload(data string, delay time.Duration) *TokenPayload {
	payload := &TokenPayload{
		ID:        uuid.NewString(),
		Data:      data,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(delay),
	}

	return payload
}

// TokenMaker 定义 jwttoken 接口两个核心接口
type TokenMaker interface {
	// GenToken 根据用户id生成有时效的token
	GenToken(*TokenPayload) (string, error)

	// ParseToken 解析并验证token是否有效
	ParseToken(token string) (*TokenPayload, error)
}

// JWTMaker jwt token 生产/解析结构体
type JWTMaker struct {
	secretKey string // 密钥
	issuer    string // 签发人
}

// NewJWTMaker 创建一个维护Token生成、解析的对象，secretKey：密钥，issuer签发token主体
func NewJWTMaker(secretKey string, issuer string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrSecretKeySize
	}
	if issuer == "" {
		return nil, ErrIssuerRequired
	}

	maker := &JWTMaker{
		secretKey: secretKey,
		issuer:    issuer,
	}

	return maker, nil
}

// GenToken 生成 Token
func (maker *JWTMaker) GenToken(payload *TokenPayload) (string, error) {
	// 声明
	var claims jwt.Claims

	// 唯一标识
	claims.ID = payload.ID
	// 主体唯一标识，可以用于存用户唯一标识，面向用户或者说使用者标识
	claims.Subject = payload.Data
	// 签发时间
	claims.Issued = jwt.NewNumericTime(payload.IssuedAt)
	// 令牌生效时间
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	// 过期时间
	claims.Expires = jwt.NewNumericTime(payload.ExpiredAt)
	// 签发人
	claims.Issuer = maker.issuer
	// Audience是指令牌的受众，通常，Audience指定为服务端的标识符
	claims.Audiences = []string{maker.issuer}

	// 用密钥签名 JWT
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(maker.secretKey))
	if err != nil {
		return "", err
	}

	return string(jwtBytes), nil
}

// ParseToken 解析并验证 Token
func (maker *JWTMaker) ParseToken(token string) (*TokenPayload, error) {
	claims, err := jwt.HMACCheck([]byte(token), []byte(maker.secretKey))
	if err != nil {
		return nil, err
	}

	if !claims.Valid(time.Now()) {
		return nil, ErrExpiredToken
	}

	if claims.Issuer != maker.issuer || !claims.AcceptAudience(maker.issuer) {
		return nil, ErrInvalidToken
	}

	payload := &TokenPayload{
		ID:        claims.ID,
		Data:      claims.Subject,
		IssuedAt:  claims.Issued.Time().Local(),
		ExpiredAt: claims.Expires.Time().Local(),
	}

	return payload, nil
}
