package gapi

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorization = "authorization"
)

func (server *Server) authorizeUser(c context.Context, username string) error {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return fmt.Errorf("metadata is not provided")
	}
	authorization := md.Get(authorization)
	if len(authorization) == 0 {
		return fmt.Errorf("authorization token is not provided")
	}
	authHeader := authorization[0]
	fields := strings.Fields(authHeader)
	if (len(fields) != 2) || strings.ToLower(fields[0]) != "bearer" {
		return fmt.Errorf("invalid authorization header format")
	}
	payload, err := server.tokenMaker.VerifyToken(fields[1])
	if err != nil {
		return err
	}
	if payload.Username != username {
		return fmt.Errorf("user is not authorized")
	}

	return nil
}
