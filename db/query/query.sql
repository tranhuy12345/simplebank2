-- name: ListQuery :many
SELECT *
FROM accounts
LEFT JOIN entries
ON accounts.id = entries.account_id;