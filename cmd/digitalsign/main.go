package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/rshep3087/coffeehouse/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger("digitalsign")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	defer func() {
		if err := log.Sync(); err != nil {
			log.Error("Error syncing log:", err)
		}
	}()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	log.Info("Starting digitalsign")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	defer nc.Close()

	log.Info("Connected to NATS")

	rnSub, err := nc.Subscribe("recipe.new", func(m *nats.Msg) {
		log.Infof("Received a message: %s", string(m.Data))
	})
	if err != nil {
		return fmt.Errorf("nats subscribe: %w", err)
	}

	defer func() {
		err := rnSub.Unsubscribe()
		if err != nil {
			log.Error("Error unsubscribing from recipe.new:", err)
		}
	}()

	// wait for ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info("Shutting down digitalsign")

	return nil
}
