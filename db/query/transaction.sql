-- name: CreateTransaction :one
INSERT INTO transaction (
  id, from_customer_accounts , to_customer_accounts , amount
) VALUES (
  $1, $2 , $3 , $4
)
RETURNING *;

-- name: ListTransaction :one
SELECT * FROM transaction
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transaction
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListTransactionWithFromID :one
SELECT * FROM transaction
WHERE from_customer_accounts = $1;

-- name: ListTransactionWithToID :one
SELECT * FROM transaction
WHERE to_customer_accounts = $1;

-- name: ListTransactionWithAmount :one
SELECT * FROM transaction
WHERE amount = $1;
