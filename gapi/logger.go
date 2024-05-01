package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	c context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	startTime := time.Now()
	result, err := handler(c, req)
	endTime := time.Now()
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	log.Info().
		Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_desc", statusCode.String()).
		Dur("duration", time.Duration(endTime.Sub(startTime).Milliseconds())).
		Msg("received request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := &ResponseRecorder{w, http.StatusOK}
		handler.ServeHTTP(rw, r)
		endTime := time.Now()

		log.Info().
			Str("protocol", "HTTP").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status_code", rw.statusCode).
			Str("status_desc", http.StatusText(rw.statusCode)).
			Dur("duration", time.Duration(endTime.Sub(startTime).Milliseconds())).
			Msg("received request")
	})
}
