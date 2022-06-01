package db

import (
	"context"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomMakeTransaction(t *testing.T) Transaction {

	FromCustomer := randomMakeCustomer(t)
	ToCustomer := randomMakeCustomer(t)

	arg := CreateTransactionParams{
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

	arg := ListTransactionsParams{
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

//func TestGetTransactionWithFromCustomerID(t *testing.T) {
//	NewTransaction := randomMakeTransaction(t)
//	SelectTransaction, err := testQueries.ListTransaction(context.Background(), NewTransaction.ID)
//
//	assert.NoError(t, err)
//	assert.NotEmpty(t, SelectTransaction)
//
//	assert.Equal(t, NewTransaction.FromCustomerAccounts, SelectTransaction.FromCustomerAccounts)
//	assert.Equal(t, NewTransaction.ToCustomerAccounts, SelectTransaction.ToCustomerAccounts)
//	assert.Equal(t, NewTransaction.Amount, SelectTransaction.Amount)
//	assert.WithinDuration(t, NewTransaction.CreatedAt, SelectTransaction.CreatedAt, time.Second)
//}

//assert.NoError(t, err)
//assert.NotEmpty(t, SelectTransaction)
//
//assert.Equal(t, NewTransaction.FromCustomerAccounts, SelectTransaction.FromCustomerAccounts)
//assert.Equal(t, NewTransaction.ToCustomerAccounts, SelectTransaction.ToCustomerAccounts)
//assert.Equal(t, NewTransaction.Amount, SelectTransaction.Amount)
//assert.WithinDuration(t, EntriesList.CreatedAt, Entries.CreatedAt, time.Second)

// func TestGetEntrieByCustomerID(t *testing.T) {
// 	Entries := randomMakeEntries(t)
// 	EntriesList, err := testQueries.ListEntriesByCustomerID(context.Background(), Entries.CustomerID)

// 	assert.NoError(t, err)
// 	for _, Entry := range EntriesList {
// 		assert.NotEmpty(t, Entry.ID)
// 	}
// }
