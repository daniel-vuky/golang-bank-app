package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	"time"

	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/daniel-vuky/golang-bank-app/val"
	"github.com/daniel-vuky/golang-bank-app/worker"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		st := status.New(codes.InvalidArgument, "invalid argument")
		br := &errdetails.BadRequest{FieldViolations: violations}
		st, _ = st.WithDetails(br)
		return nil, st.Err()

	}
	hashedPassword, hashedErr := util.HashPassword(req.GetPassword())
	if hashedErr != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", hashedErr)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			opts := []asynq.Option{
				asynq.MaxRetry(3),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDitributor.DistributeTaskSendVerifyEmail(
				c,
				&worker.PayloadSendVerifyEmail{
					Username: user.Username,
				},
				opts...,
			)
		},
	}
	user, err := server.store.CreateUserTx(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, status.Errorf(codes.Internal, "username or email already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user.User),
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (unique_violation []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: err.Error(),
		})
	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "full_name",
			Description: err.Error(),
		})
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: err.Error(),
		})
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "email",
			Description: err.Error(),
		})
	}
	return unique_violation
}
