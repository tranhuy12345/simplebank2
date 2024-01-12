-- name: CreateTransfers :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount,
  created_at
) VALUES (
  $1, $2 , $3, $4
)
RETURNING *;

-- name: GetTransfers :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id;

-- name: UpdateTransfers :exec
UPDATE transfers SET amount = $2
WHERE id = $1;

-- name: DeleteTransfers :exec
DELETE FROM transfers WHERE id = $1;