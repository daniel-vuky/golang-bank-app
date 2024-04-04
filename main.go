package main

import (
	"database/sql"
	"github.com/daniel-vuky/golang-bank-app/api"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/util"
	"log"
)

func main() {
	config, loadConfigErr := util.LoadConfig(".")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}
	conn, connectErr := sql.Open(config.DBDriver, config.DBSource)
	if connectErr != nil {
		log.Fatal("can not connect to the database", connectErr)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("fail to init the server", err)
	}

	startErr := server.Start(config.ServerAddress)
	if startErr != nil {
		log.Fatal("can not start the server", startErr)
	}
}
