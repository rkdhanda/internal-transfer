package api

import (
	"encoding/json"
	"errors"
	"internal-transfers/internal/models"
	"log/slog"
	"net/http"
)

// API error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// generate a JSON response format
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func writeError(w http.ResponseWriter, err error, defaultStatus int) {
	var status int
	var code string
	var details string

	// Map model errors to HTTP status codes
	switch {
	case errors.Is(err, models.ErrAccountNotFound):
		status = http.StatusNotFound
		code = "ACCOUNT_NOT_FOUND"
		details = err.Error()
	case errors.Is(err, models.ErrAccountExists):
		status = http.StatusConflict
		code = "ACCOUNT_EXISTS"
		details = err.Error()
	case errors.Is(err, models.ErrNegativeBalance):
		status = http.StatusBadRequest
		code = "NEGATIVE_BALANCE"
		details = err.Error()
	case errors.Is(err, models.ErrInvalidAccountID):
		status = http.StatusBadRequest
		code = "INVALID_ACCOUNT_ID"
		details = err.Error()
	case errors.Is(err, models.ErrInsufficientBalance):
		status = http.StatusUnprocessableEntity
		code = "INSUFFICIENT_BALANCE"
		details = err.Error()
	case errors.Is(err, models.ErrSelfTransfer):
		status = http.StatusBadRequest
		code = "SELF_TRANSFER"
		details = err.Error()
	case errors.Is(err, models.ErrInvalidAmount):
		status = http.StatusBadRequest
		code = "INVALID_AMOUNT"
		details = err.Error()
	case errors.Is(err, models.ErrAccountsNotFound):
		status = http.StatusBadRequest
		code = "ACCOUNTS_NOT_FOUND"
		details = err.Error()
	default:
		status = defaultStatus
		code = "INTERNAL_ERROR"
		details = "An internal error occurred"
		slog.Error("unhandled error", slog.String("error", err.Error()))
	}

	writeJSON(w, status, ErrorResponse{
		Error:   err.Error(),
		Code:    code,
		Details: details,
	})
}

func validateJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
