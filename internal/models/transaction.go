package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusPending   TransactionStatus = "pending"
)

type Transaction struct {
	TransactionID        int64             `json:"transaction_id"`
	SourceAccountID      int64             `json:"source_account_id"`
	DestinationAccountID int64             `json:"destination_account_id"`
	Amount               decimal.Decimal   `json:"amount"`
	Status               TransactionStatus `json:"status"`
	CreatedAt            time.Time         `json:"created_at"`
	ErrorMessage         *string           `json:"error_message,omitempty"`
}

type CreateTransactionRequest struct {
	SourceAccountID      int64           `json:"source_account_id" validate:"required,gt=0"`
	DestinationAccountID int64           `json:"destination_account_id" validate:"required,gt=0,nefield=SourceAccountID"`
	Amount               decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

type TransactionResponse struct {
	TransactionID        int64  `json:"transaction_id"`
	SourceAccountID      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"` // String to avoid JSON float precision issues
	Status               string `json:"status"`
}
