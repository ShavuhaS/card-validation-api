package api

import (
	"fmt"
	"net/http"
)

func NewRouter(cardHandler *CardHandler) *http.ServeMux {
	mux := &http.ServeMux{}
	
	mux.HandleFunc("POST /validate", cardHandler.HandleValidate)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Healthy")
	})

	return mux
}
