package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "github.com/freedommmoto/test_simplebank/db/mock"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeNewUser(t *testing.T) {
	userStuct, passwordString := makeNewRandomUser(t)
	fmt.Println(userStuct, "userStuct")
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(reCoder *httptest.ResponseRecorder)
	}{
		{
			name: "check case ok",
			body: gin.H{
				"username":  userStuct.Username,
				"password":  passwordString,
				"full_name": userStuct.FullName,
				"email":     userStuct.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(userStuct, nil)
			},
			checkResponse: func(reCoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, reCoder.Code)
				requireBodyMatchUser(t, reCoder.Body, userStuct)
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	dataBody, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var selcetUser db.User
	err = json.Unmarshal(dataBody, &selcetUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, selcetUser.Username)
	require.Equal(t, user.FullName, selcetUser.FullName)
	require.Equal(t, user.Email, selcetUser.Email)
	require.Empty(t, selcetUser.HashedPassword)
}

func makeNewRandomUser(t *testing.T) (User db.User, planPassword string) {
	planPassword = tool.RandomString(8)
	hashPassword, err := tool.HashPassword(planPassword)
	require.NoError(t, err)
	name := tool.RandomOwner()

	User = db.User{
		Username:       name,
		HashedPassword: hashPassword,
		FullName:       name,
		Email:          tool.RandomEmail(),
	}
	return
}
