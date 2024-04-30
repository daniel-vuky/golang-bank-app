package gapi

import (
	"fmt"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/token"
	"github.com/daniel-vuky/golang-bank-app/util"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedBankAppServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, tokenMakerErr := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if tokenMakerErr != nil {
		return nil, fmt.Errorf("can not init token maker, %w", tokenMakerErr)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
