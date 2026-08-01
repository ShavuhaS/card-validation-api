package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ShavuhaS/card-validation-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	mux := &http.ServeMux{}

	mux.HandleFunc("POST /cards/validate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	slog.Info("Starting HTTP server...", "port", cfg.Port)

	port := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start HTTP server on port %v: %v", cfg.Port, err)
	}
}
