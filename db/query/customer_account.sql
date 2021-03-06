-- name: CreateCustomer :one
INSERT INTO customer_accounts (customer_name, balance, currency)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetCustomer :one
SELECT *
FROM customer_accounts
WHERE id = $1 LIMIT 1;

-- name: GetCustomerForUpdate :one
SELECT *
FROM customer_accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListCustomer :many
SELECT *
FROM customer_accounts
ORDER BY id LIMIT $1
OFFSET $2;

-- name: ListCustomerWithOwner :many
SELECT *
FROM customer_accounts
WHERE customer_name = $1
ORDER BY id LIMIT $2
OFFSET $3;

-- name: DeleteCustomer :exec
DELETE
FROM customer_accounts
WHERE id = $1;

-- name: UpdateCustomer :one
UPDATE customer_accounts
SET balance = $2
WHERE id = $1 RETURNING *;

-- name: UpdateCustomerBalance :one
UPDATE customer_accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: UpdateCustomerCurrency :one
UPDATE customer_accounts
SET currency = $2
WHERE id = $1 RETURNING *;