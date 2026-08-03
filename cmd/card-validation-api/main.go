package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ShavuhaS/card-validation-api/internal/api"
	"github.com/ShavuhaS/card-validation-api/internal/card"
	"github.com/ShavuhaS/card-validation-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	cardService, err := card.NewValidationService()
	if err != nil {
		log.Fatalf("Failed to instantiate card validation service: %v\n", err)
	}

	cardHandler := api.NewCardHandler(cardService)

	router := api.NewRouter(cardHandler)

	slog.Info("Starting HTTP server...", "port", cfg.Port)

	port := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start HTTP server on port %v: %v", cfg.Port, err)
	}
}
