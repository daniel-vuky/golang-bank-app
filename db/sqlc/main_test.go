package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, loadConfigErr := util.LoadConfig("../..")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}
	var err error
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}
	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
