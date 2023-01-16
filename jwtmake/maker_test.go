package jwtmake_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lightsaid/gotk/jwtmake"
	"github.com/lightsaid/gotk/random"
	"github.com/stretchr/testify/require"
)

var testMaker *jwtmake.Maker

func createToken(t *testing.T) string {
	secret := random.RandomString(15)
	payload := jwtmake.NewJWTPayload("100", jwt.RegisteredClaims{})
	_, err := jwtmake.NewMaker(secret)
	require.Error(t, err)

	secret = random.RandomString(16)
	testMaker, err = jwtmake.NewMaker(secret)
	require.NoError(t, err)

	token, err := testMaker.GenToken(payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

func TestGenToken(t *testing.T) {
	createToken(t)
}

func TestParseToken(t *testing.T) {
	token := createToken(t)

	payload, err := testMaker.ParseToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, "100", payload.UID)
	require.Equal(t, jwt.NewNumericDate(time.Now().Add(15*time.Minute)), payload.ExpiresAt)
}

func TestExpiredToken(t *testing.T) {
	secret := random.RandomString(32)
	maker, err := jwtmake.NewMaker(secret)
	require.NoError(t, err)

	payload := jwtmake.NewJWTPayload("200", jwt.RegisteredClaims{
		// ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Second)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
	})

	// time.Sleep(3 * time.Second)

	token, err := maker.GenToken(payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload2, err2 := maker.ParseToken(token)
	require.ErrorIs(t, err2, jwtmake.ErrExpiredToken)
	require.Empty(t, payload2)

}

func TestErrInvalidToken(t *testing.T) {
	token := createToken(t)
	token += random.RandomString(1)
	payload, err := testMaker.ParseToken(token)
	require.ErrorIs(t, jwtmake.ErrInvalidToken, err)
	require.Empty(t, payload)
}
