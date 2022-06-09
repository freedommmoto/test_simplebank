package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/freedommmoto/test_simplebank/db/mock"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCustomer(t *testing.T) {
	customer := ranDomMakeCustomer()

	testCases := []struct {
		name          string
		customerID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			customerID: customer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(customer, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},

		{
			name:       "NOTFound",
			customerID: customer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(db.CustomerAccount{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
				//no need to check data is same or not for this case
				//requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},
		// todo add more cases here
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(customer, nil)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			//case want to test fail if customer add and select is not the same value
			//customer.Balance++
			//t.Logf("customer after mock %v", customer)

			url := fmt.Sprintf("/customer/id/%d", customer.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			//assert.Equal(t, http.StatusOK, recorder.Code)
			//requireBodyMatchAccount(t, recorder.Body, customer)
			tc.checkResponse(t, recorder)
		})
	}

}

func ranDomMakeCustomer() db.CustomerAccount {
	return db.CustomerAccount{
		ID:           tool.RandomInt(1, 1000),
		CustomerName: tool.RandomOwner(),
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
