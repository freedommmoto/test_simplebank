package token

import (
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
	"testing"
	"time"
)

var PasetoKeySize = chacha20poly1305.KeySize

func TestPasetoCreateToken(t *testing.T) {
	maker, err := NewPasetoMaker(tool.RandomString(PasetoKeySize))
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

func TestPasetoSizeNotCorrect(t *testing.T) {
	maker, err := NewPasetoMaker(tool.RandomString(PasetoKeySize + 1))
	require.Error(t, err)
	require.Empty(t, maker)
}

func TestPasetoExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(tool.RandomString(PasetoKeySize))
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
