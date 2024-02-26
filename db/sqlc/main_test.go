package db

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"database/sql"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	loadConfigErr := godotenv.Load("../../.env")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}
	var err error
	testDb, err = sql.Open("postgres", os.Getenv("POSTGRES_SOURCE"))
	if err != nil {
		log.Fatal("can not connect to database", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
