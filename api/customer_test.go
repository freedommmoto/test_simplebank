package api

import (
	"fmt"
	mockdb "github.com/freedommmoto/test_simplebank/db/mock"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCustomer(t *testing.T) {
	customer := ranDomMakeCustomer()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).Times(1).Return(customer, nil)
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	t.Logf("customer after mock %v", customer)

	url := fmt.Sprintf("/customer/id/%d", customer.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

}

func ranDomMakeCustomer() db.CustomerAccount {
	return db.CustomerAccount{
		ID:           tool.RandomInt(1, 1000),
		CustomerName: tool.RandomOwner(),
		Balance:      tool.RandomMoney(),
	}
}
