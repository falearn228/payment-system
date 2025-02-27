package api

import (
	"context"
	db "payment-system/internal/db/sqlc"
	"payment-system/internal/util"
)

func InitializeWallets(store db.Store) error {
    count, err := store.GetWalletCount(context.Background())
    if err != nil {
        return err
    }
    if count > 0 {
        return nil
    }

    for i := 0; i < 10; i++ {
        address := util.GenerateRandomAddress()
        err := store.CreateWallet(context.Background(), address)
        if err != nil {
            return err
        }
    }
    return nil
}