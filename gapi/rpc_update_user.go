package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(c context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := server.authorizeUser(c, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cannot authorize user: %v", err)
	}
	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{String: req.GetFullName(), Valid: req.FullName != nil},
		Email:    sql.NullString{String: req.GetEmail(), Valid: req.Email != nil},
	}
	if req.Password != nil {
		hashedPassword, hashedErr := util.HashPassword(req.GetPassword())
		if hashedErr != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %v", hashedErr)
		}
		arg.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}
		arg.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
	user, err := server.store.UpdateUser(c, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update user: %v", err)
	}
	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}
