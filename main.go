package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/peterbourgon/ff/v3"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/rshep3087/coffeehouse/cache/redis"
	"github.com/rshep3087/coffeehouse/database"
	"github.com/rshep3087/coffeehouse/logger"
	"github.com/rshep3087/coffeehouse/postgres"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	log, err := logger.NewLogger("coffeehouse")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		if err := log.Sync(); err != nil {
			fmt.Println("Error syncing log:", err)
		}
	}()

	if err := run(ctx, os.Args, log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, log *zap.SugaredLogger) error {
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// configuration
	fs := flag.NewFlagSet("coffeehouse", flag.ContinueOnError)
	var (
		listenAddr = fs.String("listen-addr", "localhost:8080", "listen address")
		dbUser     = fs.String("db-user", "user", "database user")
		dbPassword = fs.String("db-password", "", "database password")
		dbHost     = fs.String("db-host", "localhost:5432", "database host")
		dbName     = fs.String("db-name", "coffeehousedb", "database name")
		dbTLS      = fs.Bool("db-tls", false, "diable TLS")
		natsURL    = fs.String("nats-url", nats.DefaultURL, "nats url")
		redisURL   = fs.String("redis-url", "localhost:6379", "redis url")
	)
	if err := ff.Parse(fs, args[1:], ff.WithEnvVarPrefix("COFFEEHOUSE")); err != nil {
		return fmt.Errorf("config parse: %w", err)
	}

	// print config
	log.Infof(
		"listen addr %s, db name %s, db host %s, disable tls %t nats url %s redis url %s",
		*listenAddr,
		*dbName,
		*dbHost,
		*dbTLS,
		*natsURL,
		*redisURL,
	)

	nc, err := nats.Connect(*natsURL)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	defer nc.Close()

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// setup server
	rc := redisClient.NewClient(&redisClient.Options{Addr: *redisURL})

	s := newServer(nc, redis.New(rc))
	s.log = log

	db, err := database.Open(database.Config{
		User:       *dbUser,
		Password:   *dbPassword,
		Host:       *dbHost,
		Name:       *dbName,
		DisableTLS: *dbTLS,
	})
	if err != nil {
		return err
	}

	defer func() {
		log.Info("closing db connection")
		db.Close()
	}()

	queries := postgres.New(db)
	s.queries = queries

	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: s,
	}

	go func() {
		log.Info("server starting")
		err = srv.ListenAndServe()
		if err != nil {
			log.Errorw("server error", "ERROR", err)
			stop <- syscall.SIGTERM
		}
	}()

	<-stop

	log.Info("server stopping")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)

	if err != nil {
		log.Errorw("server shutdown error", "ERROR", err)
	}

	return nil
}
