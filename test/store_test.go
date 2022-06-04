package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMakeTransferTx(t *testing.T) {

	store := NewStore(testDB)

	customer1 := RandomMakeCustomer(t)
	customer2 := RandomMakeCustomer(t)
	fmt.Println(">> before:", customer1.Balance, customer2.Balance)

	testRound := 2
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
		require.Equal(t, amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.Amount)

		selectEntry, err := store.ListEntries(context.Background(), fromEntry.ID)
		fmt.Println(">> Entry:", selectEntry)
		require.NoError(t, err)
		require.NotEmpty(t, selectEntry)
	}

}
