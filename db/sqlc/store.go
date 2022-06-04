package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Queries
	MakeTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) SQLStore {
	return SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// make first function name as a lower case is like a private function
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb error %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transaction   Transaction     `json:"transaction"`
	FromAccountID CustomerAccount `json:"from_account_id"`
	ToAccountID   CustomerAccount `json:"to_account_id"`
	FromEntry     Entry           `json:"from_entry"`
	ToEntry       Entry           `json:"to_entry"`
}

func (store *SQLStore) MakeTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var returnData TransferTxResult
	err := store.execTx(ctx, func(queries *Queries) (err error) {

		//var err error
		returnData.Transaction, err = queries.CreateTransaction(ctx, CreateTransactionParams{
			FromCustomerAccounts: arg.FromAccountID,
			ToCustomerAccounts:   arg.ToAccountID,
			Amount:               arg.Amount,
		})
		if err != nil {
			return err
		}

		returnData.FromEntry, err = queries.CreateEntries(ctx, CreateEntriesParams{
			CustomerID: arg.FromAccountID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		returnData.ToEntry, err = queries.CreateEntries(ctx, CreateEntriesParams{
			CustomerID: arg.ToAccountID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		//todo update account amount
		
		return nil
	})
	return returnData, err
}
