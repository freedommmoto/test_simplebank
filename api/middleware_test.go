package api

import (
	"fmt"
	token "github.com/freedommmoto/test_simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addAuthHeader(t *testing.T, request *http.Request, tokerMaker token.Maker, authType string, username string, duration time.Duration) {
	tokenKey, err := tokerMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authHader := fmt.Sprintf("%s %s", authType, tokenKey)
	//
	//fmt.Println("authHader")
	//fmt.Println(username)
	//fmt.Println(authHader)

	request.Header.Set(authHeaderKeyWord, authHader)
}

func TestServer_authMiddleware(t *testing.T) {
	testCases := []struct {
		name            string
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkReturnData func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "http:ok:200",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, "testuser", time.Minute)
				//fmt.Println(request.Header)
			},
			checkReturnData: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//fmt.Println("TestServer_authMiddleware")
				//fmt.Println(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "http:noauth:401",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//case not set auth key header
			},
			checkReturnData: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "http:auth_not_support:401",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, "API Key", "testuser", time.Minute)
			},
			checkReturnData: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "http:auth_not_sent:401",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, "", "testuser", time.Minute)
			},
			checkReturnData: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "http:auth_time_out:401",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authHeaderKeyWord, "testuser", -time.Minute)
			},
			checkReturnData: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testcase := testCases[i]

		t.Run(testcase.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(
				authPath, authMiddleware(server.tokenMaker),
				func(context *gin.Context) {
					context.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testcase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testcase.checkReturnData(t, recorder)

		})
	}
}
