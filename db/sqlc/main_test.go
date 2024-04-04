package db

import (
	"github.com/daniel-vuky/golang-bank-app/util"
	"log"
	"os"
	"testing"

	"database/sql"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, loadConfigErr := util.LoadConfig("../..")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}
	var err error
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
