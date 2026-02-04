package service

import (
	"context"
	"errors"
	"fmt"
	"internal-transfers/internal/models"
	"internal-transfers/internal/repository"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type TransferService interface {
	ExecuteTransfer(ctx context.Context, req *models.CreateTransactionRequest) (*models.Transaction, error)
}

type transferService struct {
	db          *pgxpool.Pool
	accountRepo repository.AccountRepository
	txRepo      repository.TransactionRepository
	logger      *slog.Logger
}

func NewTransferService(
	db *pgxpool.Pool,
	accountRepo repository.AccountRepository,
	txRepo repository.TransactionRepository,
	logger *slog.Logger,
) TransferService {
	return &transferService{
		db:          db,
		accountRepo: accountRepo,
		txRepo:      txRepo,
		logger:      logger,
	}
}

func (s *transferService) ExecuteTransfer(ctx context.Context, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	// validate amount is positive
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		s.logger.Warn("attempted transfer with invalid amount",
			slog.String("amount", req.Amount.String()),
		)
		return nil, models.ErrInvalidAmount
	}

	// Validate source and destination are different
	if req.SourceAccountID == req.DestinationAccountID {
		s.logger.Warn("attempted self-transfer",
			slog.Int64("account_id", req.SourceAccountID),
		)
		return nil, models.ErrSelfTransfer
	}

	// Validate account IDs are positive
	if req.SourceAccountID <= 0 {
		return nil, fmt.Errorf("%w: source account_id %d", models.ErrInvalidAccountID, req.SourceAccountID)
	}
	if req.DestinationAccountID <= 0 {
		return nil, fmt.Errorf("%w: destination account_id %d", models.ErrInvalidAccountID, req.DestinationAccountID)
	}

	// Begin database transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", slog.String("error", err.Error()))
		return nil, err
	}
	defer tx.Rollback(ctx)

	firstID, secondID := req.SourceAccountID, req.DestinationAccountID
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	firstAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, firstID)
	if err != nil {
		s.logger.Error("failed to lock first account",
			slog.Int64("account_id", firstID),
			slog.String("error", err.Error()),
		)
		if errors.Is(err, models.ErrAccountNotFound) {
			return nil, fmt.Errorf("%w: account_id %d", models.ErrAccountNotFound, firstID)
		}
		return nil, err
	}

	secondAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, secondID)
	if err != nil {
		s.logger.Error("failed to lock second account",
			slog.Int64("account_id", secondID),
			slog.String("error", err.Error()),
		)
		if errors.Is(err, models.ErrAccountNotFound) {
			return nil, fmt.Errorf("%w: account_id %d", models.ErrAccountNotFound, secondID)
		}
		return nil, err
	}

	// Map accounts to source and destination
	var sourceAccount, destAccount *models.Account
	if firstID == req.SourceAccountID {
		sourceAccount = firstAccount
		destAccount = secondAccount
	} else {
		sourceAccount = secondAccount
		destAccount = firstAccount
	}

	// Check if source account has sufficient balance
	if sourceAccount.Balance.LessThan(req.Amount) {
		s.logger.Warn("insufficient balance for transfer",
			slog.Int64("source_account", req.SourceAccountID),
			slog.String("balance", sourceAccount.Balance.String()),
			slog.String("amount", req.Amount.String()),
		)

		errorMsg := "insufficient balance!"
		failedTx := &models.Transaction{
			SourceAccountID:      req.SourceAccountID,
			DestinationAccountID: req.DestinationAccountID,
			Amount:               req.Amount,
			Status:               models.TransactionStatusFailed,
			ErrorMessage:         &errorMsg,
		}
		s.txRepo.Create(ctx, tx, failedTx)
		tx.Commit(ctx)

		return nil, models.ErrInsufficientBalance
	}

	// Calculate new balances
	newSourceBalance := sourceAccount.Balance.Sub(req.Amount)
	newDestBalance := destAccount.Balance.Add(req.Amount)

	// Update balances
	err = s.accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance)
	if err != nil {
		s.logger.Error("failed to update source account balance",
			slog.Int64("account_id", req.SourceAccountID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = s.accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance)
	if err != nil {
		s.logger.Error("failed to update destination account balance",
			slog.Int64("account_id", req.DestinationAccountID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	transaction := &models.Transaction{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
		Status:               models.TransactionStatusCompleted,
	}

	transactionID, err := s.txRepo.Create(ctx, tx, transaction)
	if err != nil {
		s.logger.Error("failed to create transaction record",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	transaction.TransactionID = transactionID

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("failed to commit transaction", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("transfer completed successfully",
		slog.Int64("transaction_id", transactionID),
		slog.Int64("source_account", req.SourceAccountID),
		slog.Int64("destination_account", req.DestinationAccountID),
		slog.String("amount", req.Amount.String()),
	)

	return transaction, nil
}
