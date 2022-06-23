package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/freedommmoto/test_simplebank/db/mock"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/token"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetCustomer(t *testing.T) {
	user, _ := makeNewRandomUser(t)
	customer := ranDomMakeCustomer(user.Username)

	testCases := []struct {
		name          string
		customerID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},
		{
			name:       "NOTFound",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(db.CustomerAccount{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
				//no need to check data is same or not for this case
				//requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},
		{
			name:       "internalError",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(db.CustomerAccount{},
					sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "wrong_user_call",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, user.Username+"test", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "no_auth",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//fmt.Println(recorder.Body)
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "badRequest",
			customerID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeSupport, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(customer, nil)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//case want to test fail if customer add and select is not the same value
			//customer.Balance++
			//t.Logf("customer after mock %v", customer)

			url := fmt.Sprintf("/customer/id/%d", tc.customerID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker) // this line should start before server.router.ServeHTTP(recorder, request)
			server.router.ServeHTTP(recorder, request)
			//assert.Equal(t, http.StatusOK, recorder.Code)
			//requireBodyMatchAccount(t, recorder.Body, customer)
			tc.checkResponse(t, recorder)
		})
	}
}

func ranDomMakeCustomer(CustomerName string) db.CustomerAccount {
	return db.CustomerAccount{
		ID:           tool.RandomInt(1, 1000),
		CustomerName: CustomerName,
		Balance:      tool.RandomMoney(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, customer db.CustomerAccount) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotCustomer db.CustomerAccount
	err = json.Unmarshal(data, &gotCustomer)
	assert.NoError(t, err)
	assert.NotEmpty(t, gotCustomer)
	assert.Equal(t, customer, gotCustomer)
}
