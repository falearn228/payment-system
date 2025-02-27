package api

import (
	"net/http"
	db "payment-system/internal/db/sqlc"
)

func NewServer(store db.Store) *http.ServeMux {
	mux := http.NewServeMux()
    mux.HandleFunc("/api/send", sendHandler(store))
    mux.HandleFunc("/api/transactions", getLastHandler(store))
    mux.HandleFunc("/api/wallet/", getBalanceHandler(store))

	return mux
}