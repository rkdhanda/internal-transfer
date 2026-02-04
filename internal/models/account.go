package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountID int64           `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateAccountRequest struct {
	AccountID      int64           `json:"account_id" validate:"required,gt=0"`
	InitialBalance decimal.Decimal `json:"initial_balance" validate:"required,gte=0"`
}

type AccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"` // String to avoid JSON float precision issues
}
