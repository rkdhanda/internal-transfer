package api

import (
	"internal-transfers/internal/models"
	"internal-transfers/internal/service"
	"log/slog"
	"net/http"
)

type TransactionHandler struct {
	service service.TransferService
	logger  *slog.Logger
}

func NewTransactionHandler(service service.TransferService, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}

// handle POST /transactions
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTransactionRequest

	if err := validateJSON(r, &req); err != nil {
		h.logger.Warn("invalid JSON in create transaction request", slog.String("error", err.Error()))
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON",
			Code:  "INVALID_JSON",
		})
		return
	}

	transaction, err := h.service.ExecuteTransfer(r.Context(), &req)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := models.TransactionResponse{
		TransactionID:        transaction.TransactionID,
		SourceAccountID:      transaction.SourceAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Amount:               transaction.Amount.String(),
		Status:               string(transaction.Status),
	}

	writeJSON(w, http.StatusCreated, response)
}
