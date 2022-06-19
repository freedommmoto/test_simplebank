package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeyLong = 42

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLong {
		return nil, fmt.Errorf("key is too short need min %d", minSecretKeyLong)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	//add this for fix key is of invalid type
	//https://stackoverflow.com/questions/62229639/using-jwt-go-library-key-is-of-invalid-type
	//jwt.New(jwt.SigningMethodES512)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		varr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(varr.Inner, ErrorExporedToken) {
			return nil, ErrorExporedToken
		}
		return nil, ErrorInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}
