package gapi

import (
	"context"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, hashedErr := util.HashPassword(req.GetPassword())
	if hashedErr != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", hashedErr)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, status.Errorf(codes.Internal, "username or email already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}
	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}
