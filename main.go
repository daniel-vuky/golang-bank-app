package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/daniel-vuky/golang-bank-app/api"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	_ "github.com/daniel-vuky/golang-bank-app/doc/statik"
	"github.com/daniel-vuky/golang-bank-app/gapi"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	fs2 "github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, loadConfigErr := util.LoadConfig(".")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}
	conn, connectErr := sql.Open(config.DBDriver, config.DBSource)
	if connectErr != nil {
		log.Fatal("can not connect to the database", connectErr)
	}

	// TODO: run db migration
	runDbMigrations(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDbMigrations(url, dbSource string) {
	migration, err := migrate.New(url, dbSource)
	if err != nil {
		log.Fatal("can not create migration instance", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("can not apply migration", err)
	}

	log.Print("migration applied")
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("fail to init the server", err)
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterBankAppServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
	log.Printf("start gRPC server on %s", config.GrpcServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("fail to init the server", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankAppHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("fail to register gRPC server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFs, err := fs2.New()
	if err != nil {
		log.Fatal("can not start the server", err)
	}
	swaggerHandler := http.StripPrefix("/swagger", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
	log.Printf("start HTTP gateway server on %s", config.ServerAddress)
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("fail to init the server", err)
	}

	startErr := server.Start(config.ServerAddress)
	if startErr != nil {
		log.Fatal("can not start the server", startErr)
	}
}
