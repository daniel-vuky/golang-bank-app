package gapi

import (
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/daniel-vuky/golang-bank-app/worker"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func NewTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		PasetoSymmetricKey:  util.RandomString(32),
		PasetoTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
