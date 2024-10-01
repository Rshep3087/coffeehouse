package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nats-io/nats.go"
	"github.com/rshep3087/coffeehouse/postgres"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	defer nc.Close()

	msgChannel := make(chan postgres.Recipe)
	rnSub, err := nc.Subscribe("recipe.new", func(m *nats.Msg) {
		var r postgres.Recipe
		err := json.Unmarshal(m.Data, &r)
		if err != nil {
			return
		}

		msgChannel <- r
	})
	if err != nil {
		return fmt.Errorf("nats subscribe: %w", err)
	}

	defer func() {
		err := rnSub.Unsubscribe()
		if err != nil {
			fmt.Println("Error unsubscribing from recipe.new:", err)
		}
	}()

	p := tea.NewProgram(newModel(msgChannel))
	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("tea run: %w", err)
	}

	return nil
}
