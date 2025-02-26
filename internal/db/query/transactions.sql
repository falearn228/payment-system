-- name: GetLastTransactions :many
SELECT id, sender_address, receiver_address, amount, timestamp 
FROM transactions 
ORDER BY timestamp DESC 
LIMIT $1;

-- name: GetTransactionByID :one
SELECT * FROM transactions 
WHERE id = $1;