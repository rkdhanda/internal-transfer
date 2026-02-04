# Internal Transfers System

This system handles account creation, balance queries, and financial transactions with a focus on data integrity, consistency, and financial accuracy.


## Technology Stack

- **Language**: Go 1.22+
- **HTTP Router**: Chi 
- **Database**: PostgreSQL 16
- **Database Driver**: pgx v5
- **Migrations**: golang-migrate
- **Logging**: slog (Go standard library)
- **Containerization**: Docker & Docker Compose

## Prerequisites

- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)


## Run the application server

###  Using Docker Compose

```bash
# Clone the repository
git clone <repository-url>
cd internal-transfer

# Create a .env file based on the example
cp .env.example .env

# Start all services (PostgreSQL + API)
make docker-up

# The API will be available at http://localhost:8080
# Test the health endpoint
curl http://localhost:8080/health
```

## Configuration

  The project uses a `.env` file for configuration. A default file is created during setup, or you can create one manually:

  ```bash
  cp .env.example .env
  ```

  Ensure your `.env` contains:
  ```env
  DB_HOST=postgres or 'localhost' if running locally
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=postgres
  DB_NAME=transfers
  DB_SSL_MODE=disable
  SERVER_PORT=8080
  SERVER_HOST=0.0.0.0
  LOG_LEVEL=info
  ```




## Assumptions

1. **Single Currency**: All accounts use the same currency
2. **User-Provided IDs**: Account IDs are provided by clients (not auto-generated)
3. **Decimal Precision**: Supports up to 18 decimal places (sufficient for cryptocurrency)
4. **No Authentication**: Authentication/authorization is handled by API gateway
5. **Synchronous Processing**: Transfers are processed synchronously

