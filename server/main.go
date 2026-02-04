package main

import (
	"context"
	"fmt"
	"internal-transfers/internal/api"
	"internal-transfers/internal/config"
	"internal-transfers/internal/repository"
	"internal-transfers/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Setup logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting internal transfers service")

	// Load env configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	dbPool, err := connectDB(cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	logger.Info("connected to database successfully")

	// Initialize repos
	accountRepo := repository.NewAccountRepository(dbPool)
	transactionRepo := repository.NewTransactionRepository(dbPool)

	// Initialize services
	accountService := service.NewAccountService(accountRepo, logger)
	transferService := service.NewTransferService(dbPool, accountRepo, transactionRepo, logger)

	// Initialize API service
	accountHandler := api.NewAccountHandler(accountService, logger)
	transactionHandler := api.NewTransactionHandler(transferService, logger)

	// Setup router
	router := setupRouter(accountHandler, transactionHandler, logger)

	// Setup HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		logger.Info("server starting", slog.String("address", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("server exited")
}

func connectDB(cfg config.DatabaseConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func setupRouter(
	accountHandler *api.AccountHandler,
	transactionHandler *api.TransactionHandler,
	logger *slog.Logger,
) *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(api.RecoveryMiddleware(logger))
	router.Use(api.LoggingMiddleware(logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// API routes
	router.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.CreateAccount)
		r.Get("/{account_id}", accountHandler.GetAccount)
	})

	router.Post("/transactions", transactionHandler.CreateTransaction)

	return router
}
