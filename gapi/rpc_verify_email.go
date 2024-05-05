package gapi

import (
	"context"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(c context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		st := status.New(codes.InvalidArgument, "invalid argument")
		br := &errdetails.BadRequest{FieldViolations: violations}
		st, _ = st.WithDetails(br)
		return nil, st.Err()
	}
	result, err := server.store.VerifyEmailTx(c, db.VerifyEmailTxParams{
		EmailId: req.GetEmailId(),
		Token:   req.GetToken(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify email: %v", err)
	}
	res := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}

	return res, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (unique_violation []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "email_id",
			Description: err.Error(),
		})
	}
	if err := val.ValidateVerifyToken(req.GetToken()); err != nil {
		unique_violation = append(unique_violation, &errdetails.BadRequest_FieldViolation{
			Field:       "verify_token",
			Description: err.Error(),
		})
	}

	return unique_violation
}
