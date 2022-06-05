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

var txKey = struct{}{}

func (store *SQLStore) MakeTransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var returnData TransferTxResult
	err := store.execTx(ctx, func(queries *db.Queries) (err error) {
		txName := ctx.Value(txKey)

		fmt.Println(txName, "make transfer")
		//var err error
		returnData.Transaction, err = queries.CreateTransaction(ctx, db.CreateTransactionParams{
			FromCustomerAccounts: arg.FromAccountID,
			ToCustomerAccounts:   arg.ToAccountID,
			Amount:               arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "make CreateEntries1")
		returnData.FromEntry, err = queries.CreateEntries(ctx, db.CreateEntriesParams{
			CustomerID: arg.FromAccountID,
			Amount:     -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "make CreateEntries2")
		returnData.ToEntry, err = queries.CreateEntries(ctx, db.CreateEntriesParams{
			CustomerID: arg.ToAccountID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		//sort order update customer to make no deadlock
		if arg.FromAccountID < arg.ToAccountID {
			returnData.FromCustomer, returnData.ToCustomer, err = updateBalance(ctx, queries, txName, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			returnData.FromCustomer, returnData.ToCustomer, err = updateBalance(ctx, queries, txName, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		//case test make deadlock
		//returnData.FromCustomer, returnData.ToCustomer, err = updateBalance(ctx, queries, txName, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

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

type updateBalanceParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func updateBalance(ctx context.Context,
	queries *db.Queries,
	txName any,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (customer1 db.CustomerAccount, customer2 db.CustomerAccount, err error) {
	fmt.Println(txName, "make UpdateCustomer ", accountID1, " amount: ", amount1)
	customer1, err = queries.UpdateCustomerBalance(ctx, db.UpdateCustomerBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	fmt.Println(txName, "make UpdateCustomer ", accountID2, " amount: ", amount2)
	customer2, err = queries.UpdateCustomerBalance(context.Background(), db.UpdateCustomerBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
