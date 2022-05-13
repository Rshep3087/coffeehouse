package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/rshep3087/coffeehouse/database"
	"github.com/rshep3087/coffeehouse/logger"
	"github.com/rshep3087/coffeehouse/postgres"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger("coffeehouse")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
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
	)
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("COFFEEHOUSE")); err != nil {
		return fmt.Errorf("config parse: %w", err)
	}

	// print config
	log.Infof(
		"listen addr %s, db name %s, db host %s, disable tls %t\n",
		*listenAddr,
		*dbName,
		*dbHost,
		*dbTLS,
	)
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// setup server
	s := newServer()
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

	log.Info("server starting")
	err = http.ListenAndServe(*listenAddr, s)
	if err != nil {
		return err
	}

	return nil
}
