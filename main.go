package main

import (
	"database/sql"
	"db/api"
	db "db/db/sqlc"
	"db/db/util"
	"log"
)

// const (
// 	dbDriver = "pgx"
// 	dbSource = "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable"
// 	address  = "0.0.0.0:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot load config DB:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot load config SV:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot load config SV:", err)
	}
}
