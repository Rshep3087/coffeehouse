package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rshep3087/coffeehouse/logger"
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
	// setup server
	s := newServer()

	log.Info("server starting")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		return err
	}

	return nil
}
