package api

import (
	"errors"
	"github.com/freedommmoto/test_simplebank/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeaderKeyWord = "authorization"
	authTypeSupport   = "bearer"
	authPayloadKey    = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader(authHeaderKeyWord)
		if len(authHeader) == 0 {
			err := errors.New("auth header is not sent from client")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errerrorReturn(err))
			return
		}

		field := strings.Fields(authHeader)
		//fmt.Println("authMiddleware ====")
		//fmt.Println(field)

		//check is as type and key like "bearer tokenKeyHere"
		if len(field) < 2 {
			err := errors.New("auth key is not correct format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errerrorReturn(err))
			return
		}

		authType := strings.ToLower(field[0])
		if authType != authTypeSupport {
			err := errors.New("auth type is not correct format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errerrorReturn(err))
			return
		}

		accessToken := field[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err := errors.New("accessToken is not correct format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errerrorReturn(err))
			return
		}

		context.Set(authPayloadKey, payload)
		context.Next()
	}
}
