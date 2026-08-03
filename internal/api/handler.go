package api

import (
	"encoding/json"
	"net/http"

	"github.com/ShavuhaS/card-validation-api/internal/card"
)

type CardHandler struct {
	cardService card.ValidationService
}

func NewCardHandler(cardService card.ValidationService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

func (h *CardHandler) HandleValidate(w http.ResponseWriter, r *http.Request) {
	var req card.ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := card.NewValidationResponse(card.ErrParsingBody)
		respondJSON(w, http.StatusBadRequest, response)
		return
	}

	err := h.cardService.Validate(&req)
	response := card.NewValidationResponse(err)
	respondJSON(w, http.StatusOK, response)
}
