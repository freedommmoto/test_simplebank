package test

import (
	"context"
	"github.com/freedommmoto/test_simplebank/db/sqlc"
	tool "github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomMakeTransaction(t *testing.T) db.Transaction {

	FromCustomer := RandomMakeCustomer(t)
	ToCustomer := RandomMakeCustomer(t)

	arg := db.CreateTransactionParams{
		FromCustomerAccounts: FromCustomer.ID,
		ToCustomerAccounts:   ToCustomer.ID,
		Amount:               tool.RandomMoney(),
	}
	//t.Logf("FromCustomer %v", FromCustomer.ID)
	//t.Logf("ToCustomer %v", ToCustomer.ID)
	//t.Logf("CreateTransactionParams %v", arg)

	Transaction, err := testQueries.CreateTransaction(context.Background(), arg)

	//t.Logf("Transaction %v", Transaction)

	assert.NoError(t, err)
	assert.NotEmpty(t, Transaction)

	assert.Equal(t, arg.FromCustomerAccounts, Transaction.FromCustomerAccounts)
	assert.Equal(t, arg.ToCustomerAccounts, Transaction.ToCustomerAccounts)
	assert.Equal(t, arg.Amount, Transaction.Amount)

	assert.NotZero(t, Transaction.ID)
	assert.NotZero(t, Transaction.CreatedAt)
	return Transaction
}

// run test|debug testing
func TestCreateTransaction(t *testing.T) {
	randomMakeTransaction(t)
}

func TestGetTransactions(t *testing.T) {
	//make new two for select many
	randomMakeTransaction(t)
	randomMakeTransaction(t)

	arg := db.ListTransactionsParams{
		Limit:  2,
		Offset: 1,
	}

	SelectTransaction, err := testQueries.ListTransactions(context.Background(), arg)

	assert.NoError(t, err)
	for _, Transaction := range SelectTransaction {
		assert.NotEmpty(t, Transaction.ID)
	}
}

func TestGetTransaction(t *testing.T) {
	NewTransaction := randomMakeTransaction(t)
	SelectTransaction, err := testQueries.ListTransaction(context.Background(), NewTransaction.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, SelectTransaction)

	assert.Equal(t, NewTransaction.FromCustomerAccounts, SelectTransaction.FromCustomerAccounts)
	assert.Equal(t, NewTransaction.ToCustomerAccounts, SelectTransaction.ToCustomerAccounts)
	assert.Equal(t, NewTransaction.Amount, SelectTransaction.Amount)
	assert.WithinDuration(t, NewTransaction.CreatedAt, SelectTransaction.CreatedAt, time.Second)
}

func makeManyTransaction(t *testing.T) [2]db.CustomerAccount {
	FromCustomer := RandomMakeCustomer(t)
	ToCustomer := RandomMakeCustomer(t)
	AmountForSelect := int64(1)
	for i := 0; i < 3; i++ {
		arg := db.CreateTransactionParams{
			FromCustomerAccounts: FromCustomer.ID,
			ToCustomerAccounts:   ToCustomer.ID,
			Amount:               AmountForSelect,
		}
		_, _ = testQueries.CreateTransaction(context.Background(), arg)
	}

	//incase need Slice
	//CustomerSlice := make([]db.CustomerAccount, 0, 2)
	//CustomerSlice = append(CustomerSlice, FromCustomer)
	//CustomerSlice = append(CustomerSlice, ToCustomer)
	//return CustomerSlice

	var CustomerArray [2]db.CustomerAccount
	CustomerArray[0] = FromCustomer
	CustomerArray[1] = ToCustomer
	return CustomerArray
}

func TestListTransactionWithFromID(t *testing.T) {
	CustomerArray := makeManyTransaction(t)
	SelectTransactions, _ := testQueries.ListTransactionWithFromID(context.Background(), CustomerArray[0].ID)
	assert.NotEmpty(t, SelectTransactions)

	for _, Transaction := range SelectTransactions {
		assert.Equal(t, Transaction.FromCustomerAccounts, CustomerArray[0].ID)
		assert.NotEmpty(t, Transaction.ID)
	}
}

func TestListTransactionWithToID(t *testing.T) {
	CustomerArray := makeManyTransaction(t)
	SelectTransactions, _ := testQueries.ListTransactionWithToID(context.Background(), CustomerArray[1].ID)
	assert.NotEmpty(t, SelectTransactions)

	for _, Transaction := range SelectTransactions {
		assert.Equal(t, Transaction.ToCustomerAccounts, CustomerArray[1].ID)
		assert.NotEmpty(t, Transaction.ID)
	}
}
