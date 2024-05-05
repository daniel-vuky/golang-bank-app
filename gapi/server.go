package gapi

import (
	"fmt"

	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/token"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/daniel-vuky/golang-bank-app/worker"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedBankAppServer
	config         util.Config
	store          db.Store
	tokenMaker     token.Maker
	taskDitributor worker.TaskDistributor
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, tokenMakerErr := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if tokenMakerErr != nil {
		return nil, fmt.Errorf("can not init token maker, %w", tokenMakerErr)
	}
	server := &Server{
		config:         config,
		store:          store,
		tokenMaker:     tokenMaker,
		taskDitributor: taskDistributor,
	}

	return server, nil
}
