package service

import (
	"context"
	"internal-transfers/internal/models"
	"internal-transfers/internal/repository"
	"log/slog"

	"github.com/shopspring/decimal"
)

// interface for account business logic
type AccountService interface {
	CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error)
	GetAccount(ctx context.Context, accountID int64) (*models.Account, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
	logger      *slog.Logger
}

// create a new account service
func NewAccountService(accountRepo repository.AccountRepository, logger *slog.Logger) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

// create a new account with the given initial balance
func (s *accountService) CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error) {
	// validate initial balance is not negative
	if req.InitialBalance.LessThan(decimal.Zero) {
		s.logger.Warn("attempted to create account with negative balance",
			slog.Int64("account_id", req.AccountID),
			slog.String("balance", req.InitialBalance.String()),
		)
		return nil, models.ErrNegativeBalance
	}

	// Validate account ID is positive
	if req.AccountID <= 0 {
		s.logger.Warn("attempted to create account with invalid ID",
			slog.Int64("account_id", req.AccountID),
		)
		return nil, models.ErrInvalidAccountID
	}

	account := &models.Account{
		AccountID: req.AccountID,
		Balance:   req.InitialBalance,
	}

	err := s.accountRepo.Create(ctx, account)
	if err != nil {
		s.logger.Error("failed to create account",
			slog.Int64("account_id", req.AccountID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	s.logger.Info("account created successfully",
		slog.Int64("account_id", account.AccountID),
		slog.String("balance", account.Balance.String()),
	)

	return account, nil
}

// get an account by ID
func (s *accountService) GetAccount(ctx context.Context, accountID int64) (*models.Account, error) {
	if accountID <= 0 {
		return nil, models.ErrInvalidAccountID
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("failed to get account",
			slog.Int64("account_id", accountID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return account, nil
}
