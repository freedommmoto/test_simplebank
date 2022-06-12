package test

import (
	"context"
	"database/sql"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"testing"
	"time"

	tool "github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/assert"
)

func RandomMakeCustomer(t *testing.T) db.CustomerAccount {
	user := RandomMakeUser(t)

	arg := db.CreateCustomerParams{
		CustomerName: user.Username,
		Balance:      tool.RandomMoney(),
		Currency:     tool.RandomCurrency(),
	}

	customer, err := testQueries.CreateCustomer(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, customer)

	assert.Equal(t, arg.CustomerName, customer.CustomerName)
	assert.Equal(t, arg.Balance, customer.Balance)
	assert.Equal(t, arg.Currency, customer.Currency)

	assert.NotZero(t, customer.ID)
	assert.NotZero(t, customer.CreatedAt)
	return customer
}

// run test|debug testing
func TestCreateCustomer(t *testing.T) {
	RandomMakeCustomer(t)
}

func TestGetCustomer(t *testing.T) {
	customer1 := RandomMakeCustomer(t)
	customer2, err := testQueries.GetCustomer(context.Background(), customer1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, customer2)

	assert.Equal(t, customer1.CustomerName, customer2.CustomerName)
	assert.Equal(t, customer1.Balance, customer2.Balance)
	assert.Equal(t, customer1.Currency, customer2.Currency)
	assert.WithinDuration(t, customer1.CreatedAt, customer2.CreatedAt, time.Second)
}

func TestUpdateCustomer(t *testing.T) {
	customer1 := RandomMakeCustomer(t)

	arg := db.UpdateCustomerParams{
		ID:      customer1.ID,
		Balance: tool.RandomMoney(),
	}

	customer2, err := testQueries.UpdateCustomer(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, customer2)
	assert.Equal(t, customer1.CustomerName, customer2.CustomerName)
	assert.Equal(t, arg.Balance, customer2.Balance)
	assert.Equal(t, customer1.Currency, customer2.Currency)
	assert.WithinDuration(t, customer1.CreatedAt, customer2.CreatedAt, time.Second)
}

func TestDeleteCustomer(t *testing.T) {
	customerTestDelete := RandomMakeCustomer(t)
	err := testQueries.DeleteCustomer(context.Background(), customerTestDelete.ID)
	assert.NoError(t, err)

	customer, err := testQueries.GetCustomer(context.Background(), customerTestDelete.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, customer)
}

func TestListCustomer(t *testing.T) {
	for i := 0; i < 10; i++ {
		RandomMakeCustomer(t)
	}

	arg := db.ListCustomerParams{
		Limit:  4,
		Offset: 6,
	}

	customerList, err := testQueries.ListCustomer(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, customerList, 4)

	for _, customer := range customerList {
		assert.NotEmpty(t, customer.ID)
	}

}
