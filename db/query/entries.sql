-- name: CreateEntries :one
INSERT INTO entries (
  id, customer_id , amount
) VALUES (
  $1, $2 , $3
)
RETURNING *;

-- name: ListEntries :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntriesByCustomerID :many
SELECT * FROM entries
WHERE customer_id = $1;
