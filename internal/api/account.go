package api

import (
	"internal-transfers/internal/models"
	"internal-transfers/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	service service.AccountService
	logger  *slog.Logger
}

func NewAccountHandler(service service.AccountService, logger *slog.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest

	if err := validateJSON(r, &req); err != nil {
		h.logger.Warn("invalid JSON in create account request", slog.String("error", err.Error()))
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON",
			Code:  "INVALID_JSON",
		})
		return
	}

	account, err := h.service.CreateAccount(r.Context(), &req)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	// Return account response
	response := models.AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance.String(),
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "account_id")

	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		h.logger.Warn("invalid account ID format", slog.String("account_id", accountIDStr))
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid account ID format",
			Code:  "INVALID_ID_FORMAT",
		})
		return
	}

	account, err := h.service.GetAccount(r.Context(), accountID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := models.AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance.String(),
	}

	writeJSON(w, http.StatusOK, response)
}
