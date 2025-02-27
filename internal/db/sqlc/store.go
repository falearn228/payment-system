package db

import (
	"context"
	"fmt"
	"payment-system/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, from, to string, amount decimal.Decimal) error
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// TransferTxParams для перевода между пользователями
type TransferTxParams struct {
	FromUserID int32           `json:"sender_address"`
	ToUserID   int32           `json:"receiver_address"`
	Amount     decimal.Decimal `json:"amount"`
}

// execTx обрабатывает транзакции
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := New(tx)
	err = fn(qtx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// TransferTx выполняет перевод между двумя кошельками
func (store *SQLStore) TransferTx(ctx context.Context, from, to string, amount decimal.Decimal) error {
    return store.execTx(ctx, func(q *Queries) error {
        // Конвертируем amount в pgtype.Numeric
        amountNumeric, err := util.DecimalToNumeric(amount)
        if err != nil {
            return fmt.Errorf("failed to convert amount into numeric: %v", err)
        }

        // Проверка существования и баланса отправителя с блокировкой
        balanceNumeric, err := q.GetWalletForUpdate(ctx, from)
        if err != nil {
            return fmt.Errorf("sender wallet not found: %v", err)
        }
        balance, err := util.NumericToDecimal(balanceNumeric)
        if err != nil {
            return fmt.Errorf("failed to convert sender balance: %v", err)
        }
        if balance.LessThan(amount) {
            return fmt.Errorf("insufficient funds")
        }

        // Проверка существования кошелька получателя с блокировкой
        _, err = q.GetWalletForUpdate(ctx, to)
        if err != nil {
            return fmt.Errorf("receiver wallet not found: %v", err)
        }

        // Уменьшение баланса отправителя
        err = q.DecrementBalance(ctx, DecrementBalanceParams{
            Balance: amountNumeric,
            Address: from,
        })
        if err != nil {
            return err
        }

        // Увеличение баланса получателя
        err = q.IncrementBalance(ctx, IncrementBalanceParams{
            Balance: amountNumeric,
            Address: to,
        })
        if err != nil {
            return err
        }

        // Запись транзакции
        err = q.CreateTransaction(ctx, CreateTransactionParams{
            SenderAddress:   from,
            ReceiverAddress: to,
            Amount:          amountNumeric,
        })
        if err != nil {
            return err
        }

        return nil
    })
}