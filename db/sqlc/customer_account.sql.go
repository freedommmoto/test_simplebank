// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: customer_account.sql

package db

import (
	"context"
)

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customer_accounts (customer_name, balance, currency)
VALUES ($1, $2, $3) RETURNING id, customer_name, balance, currency, created_at
`

type CreateCustomerParams struct {
	CustomerName string `json:"customer_name"`
	Balance      int64  `json:"balance"`
	Currency     string `json:"currency"`
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, createCustomer, arg.CustomerName, arg.Balance, arg.Currency)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteCustomer = `-- name: DeleteCustomer :exec
DELETE
FROM customer_accounts
WHERE id = $1
`

func (q *Queries) DeleteCustomer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCustomer, id)
	return err
}

const getCustomer = `-- name: GetCustomer :one
SELECT id, customer_name, balance, currency, created_at
FROM customer_accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCustomer(ctx context.Context, id int64) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, getCustomer, id)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getCustomerForUpdate = `-- name: GetCustomerForUpdate :one
SELECT id, customer_name, balance, currency, created_at
FROM customer_accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetCustomerForUpdate(ctx context.Context, id int64) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, getCustomerForUpdate, id)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listCustomer = `-- name: ListCustomer :many
SELECT id, customer_name, balance, currency, created_at
FROM customer_accounts
ORDER BY id LIMIT $1
OFFSET $2
`

type ListCustomerParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCustomer(ctx context.Context, arg ListCustomerParams) ([]CustomerAccount, error) {
	rows, err := q.db.QueryContext(ctx, listCustomer, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CustomerAccount
	for rows.Next() {
		var i CustomerAccount
		if err := rows.Scan(
			&i.ID,
			&i.CustomerName,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCustomerWithOwner = `-- name: ListCustomerWithOwner :many
SELECT id, customer_name, balance, currency, created_at
FROM customer_accounts
WHERE customer_name = $1
ORDER BY id LIMIT $2
OFFSET $3
`

type ListCustomerWithOwnerParams struct {
	CustomerName string `json:"customer_name"`
	Limit        int32  `json:"limit"`
	Offset       int32  `json:"offset"`
}

func (q *Queries) ListCustomerWithOwner(ctx context.Context, arg ListCustomerWithOwnerParams) ([]CustomerAccount, error) {
	rows, err := q.db.QueryContext(ctx, listCustomerWithOwner, arg.CustomerName, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CustomerAccount
	for rows.Next() {
		var i CustomerAccount
		if err := rows.Scan(
			&i.ID,
			&i.CustomerName,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCustomer = `-- name: UpdateCustomer :one
UPDATE customer_accounts
SET balance = $2
WHERE id = $1 RETURNING id, customer_name, balance, currency, created_at
`

type UpdateCustomerParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, updateCustomer, arg.ID, arg.Balance)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateCustomerBalance = `-- name: UpdateCustomerBalance :one
UPDATE customer_accounts
SET balance = balance + $1
WHERE id = $2 RETURNING id, customer_name, balance, currency, created_at
`

type UpdateCustomerBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) UpdateCustomerBalance(ctx context.Context, arg UpdateCustomerBalanceParams) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, updateCustomerBalance, arg.Amount, arg.ID)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateCustomerCurrency = `-- name: UpdateCustomerCurrency :one
UPDATE customer_accounts
SET currency = $2
WHERE id = $1 RETURNING id, customer_name, balance, currency, created_at
`

type UpdateCustomerCurrencyParams struct {
	ID       int64  `json:"id"`
	Currency string `json:"currency"`
}

func (q *Queries) UpdateCustomerCurrency(ctx context.Context, arg UpdateCustomerCurrencyParams) (CustomerAccount, error) {
	row := q.db.QueryRowContext(ctx, updateCustomerCurrency, arg.ID, arg.Currency)
	var i CustomerAccount
	err := row.Scan(
		&i.ID,
		&i.CustomerName,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
