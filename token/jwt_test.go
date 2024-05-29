package token

import (
	"db/db/util"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	var string1 = util.RandomString(32)
	fmt.Println(len(string1))
	maker, err := NewPasetoMaker(string1)
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	now := time.Now()
	issueAt := now.Unix()
	exprireAt := now.Add(duration).Unix()

	token, _, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, issueAt, payload.IssuedAt)
	require.Equal(t, exprireAt, payload.ExpiredAt)
}

func TestExpriredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, _, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	fmt.Println(err)
	require.Nil(t, payload)
}

func TestInvalidAlgNone(t *testing.T) {
	payload, err := NewPayLoad(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id":       payload.ID,
		"username": payload.Username,
		"exp":      payload.ExpiredAt,
		"issue":    payload.IssuedAt,
	})
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	fmt.Println(err)
	require.Error(t, err)
	require.Nil(t, payload)
}
