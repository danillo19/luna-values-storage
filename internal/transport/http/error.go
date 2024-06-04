package http

import (
	"encoding/json"
	"errors"
	"log"
	"luna-values-storage/internal/domain"
	"net/http"
)

var (
	ErrParsingFailed = errors.New("failed to parse request body")
)

func RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	if errors.Is(err, domain.ErrNotFound) {
		httpRespondWithError(err, "Not found", w, r, http.StatusNotFound)
	} else if errors.Is(err, domain.ErrInvalidVariable) {
		httpRespondWithError(err, "Invalid variable", w, r, http.StatusBadRequest)
	} else if errors.Is(err, domain.ErrInvalidValue) {
		httpRespondWithError(err, "Invalid value", w, r, http.StatusBadRequest)
	} else if errors.Is(err, domain.ErrRequired) {
		httpRespondWithError(err, err.Error(), w, r, http.StatusBadRequest)
	} else if errors.Is(err, ErrParsingFailed) {
		httpRespondWithError(err, err.Error(), w, r, http.StatusBadRequest)
	} else {
		httpRespondWithError(err, "internal-server-error", w, r, http.StatusInternalServerError)
	}
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, status int) {
	log.Printf("error: %s, slug: %s", err, slug)

	resp := ErrorResponse{slug, status}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
