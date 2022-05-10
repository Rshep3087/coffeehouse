package main

import (
	"log"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// setup server
	s := newServer()

	err := http.ListenAndServe(":8080", s)
	if err != nil {
		return err
	}

	return nil
}
