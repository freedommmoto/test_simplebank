-- name: CreateCustomer :one
INSERT INTO customer_accounts (customer_name, balance, currency)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetCustomer :one
SELECT *
FROM customer_accounts
WHERE id = $1 LIMIT 1;

-- name: ListCustomer :many
SELECT *
FROM customer_accounts
ORDER BY id LIMIT $1
OFFSET $2;

-- name: DeleteCustomer :exec
DELETE
FROM customer_accounts
WHERE id = $1;

-- name: UpdateCustomer :one
UPDATE customer_accounts
SET balance = $2
WHERE id = $1 RETURNING *;

-- name: UpdateCustomerCurrency :one
UPDATE customer_accounts
SET currency = $2
WHERE id = $1 RETURNING *;