package test

import (
	"context"
	"fmt"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMakeTransferTx(t *testing.T) {
	store := NewStore(testDB)

	customer1 := RandomMakeCustomer(t)
	customer2 := RandomMakeCustomer(t)

	_, err := store.UpdateCustomer(context.Background(), db.UpdateCustomerParams{
		ID:      customer1.ID,
		Balance: 100,
	})
	require.NoError(t, err)
	_, err = store.UpdateCustomer(context.Background(), db.UpdateCustomerParams{
		ID:      customer2.ID,
		Balance: 0,
	})
	require.NoError(t, err)

	customer1, err = store.GetCustomer(context.Background(), customer1.ID)
	customer2, err = store.GetCustomer(context.Background(), customer2.ID)

	fmt.Println(">> before:", customer1.Balance, customer2.Balance)

	testRound := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < testRound; i++ {
		go func() {
			result, err := store.MakeTransferTx(context.Background(), TransferTxParams{
				FromAccountID: customer1.ID,
				ToAccountID:   customer2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	//existed := make(map[int]bool)
	//check results
	for m := 0; m < testRound; m++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		Transaction := result.Transaction
		require.NotEmpty(t, Transaction)
		require.Equal(t, customer1.ID, Transaction.FromCustomerAccounts)
		require.Equal(t, customer2.ID, Transaction.ToCustomerAccounts)
		require.Equal(t, amount, Transaction.Amount)
		require.NotZero(t, Transaction.ID)
		require.NotEmpty(t, Transaction.CreatedAt)

		//t.ListTransactionWithToID(context.Background(), Transaction.ID)
		_, err = store.ListTransactionWithToID(context.Background(), Transaction.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, customer1.ID, fromEntry.CustomerID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.Amount)

		selectEntry, err := store.ListEntries(context.Background(), fromEntry.ID)
		//fmt.Println(">> Entry:", selectEntry)
		require.NoError(t, err)
		require.NotEmpty(t, selectEntry)

		// check customer
		fromCustomer := result.FromCustomer
		require.NotEmpty(t, fromCustomer)
		require.Equal(t, fromCustomer.ID, customer1.ID)

		toCustomer := result.ToCustomer
		require.NotEmpty(t, toCustomer)
		require.Equal(t, toCustomer.ID, customer2.ID)

		//check customer balance
		diff1 := customer1.Balance - fromCustomer.Balance
		diff2 := toCustomer.Balance - customer2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		fmt.Println(">> tx  :", fromCustomer.Balance, toCustomer.Balance)

		//need to check more here
		//k := int(diff1 / amount)
		//require.True(t, k >= 1 && k <= testRound)
		//require.NotContains(t, existed, k)
		//existed[k] = true
	}

	//after for loop check the final update balances
	updateCustomer1, err := testQueries.GetCustomer(context.Background(), customer1.ID)
	require.NoError(t, err)

	updateCustomer2, err := testQueries.GetCustomer(context.Background(), customer2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateCustomer1.Balance, updateCustomer2.Balance)

	require.Equal(t, customer1.Balance-(int64(testRound)*amount), updateCustomer1.Balance)
	require.Equal(t, customer1.Balance+(int64(testRound)*amount), updateCustomer2.Balance)

}
