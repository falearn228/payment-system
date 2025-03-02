// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: transactions.sql

package db

import (
	"context"
)

const getLastTransactions = `-- name: GetLastTransactions :many
SELECT id, sender_address, receiver_address, amount, timestamp
FROM transactions
ORDER BY timestamp DESC
LIMIT $1
`

func (q *Queries) GetLastTransactions(ctx context.Context, limit int32) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, getLastTransactions, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.SenderAddress,
			&i.ReceiverAddress,
			&i.Amount,
			&i.Timestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT id, sender_address, receiver_address, amount, timestamp FROM transactions 
WHERE id = $1
`

func (q *Queries) GetTransactionByID(ctx context.Context, id int32) (Transaction, error) {
	row := q.db.QueryRow(ctx, getTransactionByID, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.SenderAddress,
		&i.ReceiverAddress,
		&i.Amount,
		&i.Timestamp,
	)
	return i, err
}
