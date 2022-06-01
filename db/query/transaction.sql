-- name: CreateTransaction :one
INSERT INTO transaction (from_customer_accounts, to_customer_accounts, amount)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListTransaction :one
SELECT *
FROM transaction
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT *
FROM transaction
ORDER BY id LIMIT $1
OFFSET $2;

-- name: ListTransactionWithFromID :many
SELECT *
FROM transaction
WHERE from_customer_accounts = $1;

-- name: ListTransactionWithToID :many
SELECT *
FROM transaction
WHERE to_customer_accounts = $1;

-- name: ListTransactionWithAmount :many
SELECT *
FROM transaction
WHERE amount = $1;
