package test

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
)

type Store interface {
	db.Queries
	MakeTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*db.Queries
}

func NewStore(sqldb *sql.DB) SQLStore {
	return SQLStore{
		db:      sqldb,
		Queries: db.New(sqldb),
	}
}

// make first function name as a lower case is like a private function
func (store *SQLStore) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
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
	Transaction  db.Transaction     `json:"transaction"`
	FromCustomer db.CustomerAccount `json:"from_account_id"`
	ToCustomer   db.CustomerAccount `json:"to_account_id"`
	FromEntry    db.Entry           `json:"from_entry"`
	ToEntry      db.Entry           `json:"to_entry"`
}

func (store *SQLStore) MakeTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var returnData TransferTxResult
	err := store.execTx(ctx, func(queries *db.Queries) (err error) {

		//var err error
		returnData.Transaction, err = queries.CreateTransaction(ctx, db.CreateTransactionParams{
			FromCustomerAccounts: arg.FromAccountID,
			ToCustomerAccounts:   arg.ToAccountID,
			Amount:               arg.Amount,
		})
		if err != nil {
			return err
		}

		returnData.FromEntry, err = queries.CreateEntries(ctx, db.CreateEntriesParams{
			CustomerID: arg.FromAccountID,
			Amount:     -arg.Amount,
		})
		if err != nil {
			return err
		}

		returnData.ToEntry, err = queries.CreateEntries(ctx, db.CreateEntriesParams{
			CustomerID: arg.ToAccountID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		//select customer after update Balance
		returnData.FromCustomer, err = queries.GetCustomer(ctx, arg.FromAccountID)
		returnData.ToCustomer, err = queries.GetCustomer(ctx, arg.ToAccountID)

		//todo update account amount
		_, err = queries.UpdateCustomer(context.Background(), db.UpdateCustomerParams{
			ID:      arg.FromAccountID,
			Balance: returnData.FromCustomer.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		_, err = queries.UpdateCustomer(context.Background(), db.UpdateCustomerParams{
			ID:      arg.ToAccountID,
			Balance: returnData.ToCustomer.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		//select customer after update Balance
		returnData.FromCustomer, err = queries.GetCustomer(ctx, arg.FromAccountID)
		returnData.ToCustomer, err = queries.GetCustomer(ctx, arg.ToAccountID)

		return nil
	})
	return returnData, err
}
