package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"payment-system/internal/api"
	db "payment-system/internal/db/sqlc"
	util "payment-system/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load configurations: ", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalln("can't connect to DB: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)

	if err := api.InitializeWallets(store); err != nil {
        log.Fatal(err)
    }

    // Настройка маршрутов
    mux := api.NewServer(store);

    fmt.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
