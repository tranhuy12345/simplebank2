-- name: CreateEntries :one
INSERT INTO entries (
  account_id,
  amount,
  created_at
) VALUES (
  $1, $2 , $3
)
RETURNING *;

-- name: GetEntries :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id;

-- name: UpdateEntries :exec
UPDATE entries SET amount = $2
WHERE id = $1;

-- name: DeleteEntries :exec
DELETE FROM entries WHERE id = $1;

-- name: DeleteEntriesByAccountId :exec
DELETE FROM entries WHERE account_id = $1;


