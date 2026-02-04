package repository

import (
	"context"
	"internal-transfers/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) (int64, error)
	GetByID(ctx context.Context, transactionID int64) (*models.Transaction, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

// create a new transaction record
func (r *transactionRepository) Create(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) (int64, error) {
	query := `
		INSERT INTO transactions (
			source_account_id, 
			destination_account_id, 
			amount, 
			status, 
			created_at,
			error_message
		)
		VALUES ($1, $2, $3, $4, NOW(), $5)
		RETURNING transaction_id
	`

	var transactionID int64
	err := tx.QueryRow(
		ctx,
		query,
		transaction.SourceAccountID,
		transaction.DestinationAccountID,
		transaction.Amount,
		transaction.Status,
		transaction.ErrorMessage,
	).Scan(&transactionID)

	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

// get a transaction by ID
func (r *transactionRepository) GetByID(ctx context.Context, transactionID int64) (*models.Transaction, error) {
	query := `
		SELECT 
			transaction_id,
			source_account_id,
			destination_account_id,
			amount,
			status,
			created_at,
			error_message
		FROM transactions
		WHERE transaction_id = $1
	`

	var transaction models.Transaction
	err := r.db.QueryRow(ctx, query, transactionID).Scan(
		&transaction.TransactionID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.ErrorMessage,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
