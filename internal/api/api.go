package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	db "payment-system/internal/db/sqlc"
	"payment-system/internal/util"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type Transaction struct {
    ID              int32           `json:"id"`
    SenderAddress   string          `json:"sender_address"`
    ReceiverAddress string          `json:"receiver_address"`
    Amount          pgtype.Numeric  `json:"amount"`
    Timestamp       time.Time       `json:"timestamp"`
}


//POST /api/send
func sendHandler(store db.Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Чтение тела запроса
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        // Парсинг JSON
        var params struct {
            From   string          `json:"from"`
            To     string          `json:"to"`
            Amount decimal.Decimal `json:"amount"`
        }
        if err := json.Unmarshal(body, &params); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        // Валидация входных данных
        if params.From == "" || params.To == "" || params.Amount.LessThanOrEqual(decimal.Zero) {
            http.Error(w, "Invalid parameters", http.StatusBadRequest)
            return
        }

        // Выполнение перевода
        err = store.TransferTx(r.Context(), params.From, params.To, params.Amount)
        if err != nil {
            switch err.Error() {
            case "insufficient funds":
                http.Error(w, "Insufficient funds", http.StatusBadRequest)
            case "sender wallet not found", "receiver wallet not found":
                http.Error(w, "Wallet not found", http.StatusNotFound)
            default:
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Transaction successful"))
    }
}

//GET /api/transactions/count=N
func getLastHandler(store db.Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        nStr := r.URL.Query().Get("count")
        if nStr == "" {
            http.Error(w, "Missing count parameter", http.StatusBadRequest)
            return
        }
        n, err := strconv.ParseInt(nStr, 10, 32)
        if err != nil || n <= 0 {
            http.Error(w, "Invalid count parameter", http.StatusBadRequest)
            return
        }

        transactions, err := store.GetLastTransactions(r.Context(), int32(n))
        if err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        // Конвертация amount в decimal.Decimal для JSON
        type jsonTransaction struct {
            ID              int32           `json:"id"`
            SenderAddress   string          `json:"sender_address"`
            ReceiverAddress string          `json:"receiver_address"`
            Amount          decimal.Decimal `json:"amount"`
            Timestamp       time.Time       `json:"timestamp"`
        }
        var jsonTxs []jsonTransaction
        for _, t := range transactions {
            amount, _ := util.NumericToDecimal(t.Amount)
            jsonTxs = append(jsonTxs, jsonTransaction{
                ID:              t.ID,
                SenderAddress:   t.SenderAddress,
                ReceiverAddress: t.ReceiverAddress,
                Amount:          amount,
            })
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(jsonTxs); err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }
}

// GET /api/wallet/{address}/balance
func getBalanceHandler(store db.Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        parts := strings.Split(r.URL.Path, "/")
        if len(parts) < 4 || parts[1] != "api" || parts[2] != "wallet" || parts[4] != "balance" {
            http.Error(w, "Invalid path", http.StatusBadRequest)
            return
        }
        address := parts[3]

        balance, err := store.GetBalance(r.Context(), address)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Wallet not found", http.StatusNotFound)
            } else {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
            return
        }

        balanceDec, err := util.NumericToDecimal(balance)
        if err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }

        response := struct {
            Address string          `json:"address"`
            Balance decimal.Decimal `json:"balance"`
        }{Address: address, Balance: balanceDec }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }
}