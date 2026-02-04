package models

import "errors"

var (
	// Account errors
	ErrAccountNotFound  = errors.New("account not found")
	ErrAccountExists    = errors.New("account already exists")
	ErrNegativeBalance  = errors.New("balance cannot be negative")
	ErrInvalidAccountID = errors.New("invalid account ID")

	// Transaction errors
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSelfTransfer        = errors.New("cannot transfer to same account")
	ErrInvalidAmount       = errors.New("invalid transaction amount")
	ErrAccountsNotFound    = errors.New("one or both accounts not found")
	ErrTransactionFailed   = errors.New("transaction failed")
)
