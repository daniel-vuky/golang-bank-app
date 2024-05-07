package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/daniel-vuky/golang-bank-app/mail"

	"github.com/daniel-vuky/golang-bank-app/api"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	_ "github.com/daniel-vuky/golang-bank-app/doc/statik"
	"github.com/daniel-vuky/golang-bank-app/gapi"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/daniel-vuky/golang-bank-app/worker"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	fs2 "github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interrruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, loadConfigErr := util.LoadConfig(".")
	if loadConfigErr != nil {
		log.Fatal("can not load the config file", loadConfigErr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interrruptSignals...)
	defer stop()

	connPool, connectErr := pgxpool.New(ctx, config.DBSource)
	if connectErr != nil {
		log.Fatal("can not connect to the database", connectErr)
	}

	// TODO: run db migration
	runDbMigrations(config.MigrationUrl, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err := waitGroup.Wait()
	if err != nil {
		log.Fatal("err from wait group", err)
	}
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

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
) {
	mailer := mail.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)
	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Print("start task processor")
	err := processor.Start()
	if err != nil {
		log.Fatal("can not start the task processor", err)
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Print("shutting down task processor")
		processor.Shutdown()
		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, store, taskDistributor)
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
	waitGroup.Go(func() error {
		log.Printf("start gRPC server on %s", config.GrpcServerAddress)
		err = grpcServer.Serve(listener)
		if err != nil {
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Print("shutting down gRPC server")
		grpcServer.GracefulStop()
		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, store, taskDistributor)
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

	httpServer := &http.Server{
		Addr:    config.ServerAddress,
		Handler: gapi.HttpLogger(mux),
	}

	waitGroup.Go(func() error {
		log.Printf("start HTTP gateway server on %s", config.ServerAddress)
		err = httpServer.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Print("shutting down HTTP gateway server")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
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
