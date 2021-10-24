package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rdforte/go-service/app/handlers"
)

func main() {

	logger := log.New(os.Stdout, "Ryans Service ", log.LstdFlags)

	app := handlers.CreateApp(logger)

	s := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	if err := s.ListenAndServe(); err != nil {
		logger.Fatal("error starting server", err)
	}
}
