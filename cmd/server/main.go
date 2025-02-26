package main

import (
	"context"
	"log"
	db "payment-system/internal/db/sqlc"
	"payment-system/internal/util"

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
	server, err := api.NewServer(store)

	err = server.Start(config.Address)
	if err != nil {
		log.Fatalln("cant't start a server: ", err)
	}
}
