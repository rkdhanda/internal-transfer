package repository

import (
	"context"
	"errors"
	"internal-transfers/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, accountID int64) (*models.Account, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, accountID int64, newBalance decimal.Decimal) error
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*models.Account, error)
}

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepository{db: db}
}

// create a new account
func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO accounts (account_id, balance, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := r.db.Exec(ctx, query, account.AccountID, account.Balance)
	if err != nil {
		// Check for unique constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return models.ErrAccountExists
			}
		}
		return err
	}

	return nil
}

// get an account by ID
func (r *accountRepository) GetByID(ctx context.Context, accountID int64) (*models.Account, error) {
	query := `
		SELECT account_id, balance, created_at, updated_at
		FROM accounts
		WHERE account_id = $1
	`

	var account models.Account
	err := r.db.QueryRow(ctx, query, accountID).Scan(
		&account.AccountID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

// get an account by ID
func (r *accountRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*models.Account, error) {
	query := `
		SELECT account_id, balance, created_at, updated_at
		FROM accounts
		WHERE account_id = $1
		FOR UPDATE
	`

	var account models.Account
	err := tx.QueryRow(ctx, query, accountID).Scan(
		&account.AccountID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

// update the balance
func (r *accountRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, accountID int64, newBalance decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = NOW()
		WHERE account_id = $2
	`

	result, err := tx.Exec(ctx, query, newBalance, accountID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrAccountNotFound
	}

	return nil
}
