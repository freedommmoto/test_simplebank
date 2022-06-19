package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTCreateToken(t *testing.T) {
	maker, err := NewJWTMaker(tool.RandomString(42))
	require.NoError(t, err)

	userName := tool.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userName, payload.Username)
	require.WithinDurationf(t, issuedAt, payload.IssuedAt, time.Second, "error message %s", "formatted")
	require.WithinDurationf(t, expiredAt, payload.ExpiredAt, time.Second, "error message %s", "formatted")
}

func TestJWTExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(tool.RandomString(42))
	require.NoError(t, err)

	userName := tool.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(userName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExporedToken.Error())
	require.Nil(t, payload)
}

func TestCaseNoneJwtTokenAlgorithmType(t *testing.T) {
	maker, err := NewJWTMaker(tool.RandomString(42))
	require.NoError(t, err)

	payload, err := NewPayload(tool.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SigningString()
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrorInvalidToken.Error())
}
