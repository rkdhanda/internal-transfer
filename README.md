# Internal Transfers Service

A Go-based microservice for handling internal money transfers between accounts.

## Prerequisites

- [Go](https://go.dev/dl/) (v1.25.6 or later)
- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional, for easy commands)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

## Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/internal-transfers.git
    cd internal-transfers
    ```

2.  **Environment Variables:**

    The project uses a `.env` file for configuration. A default file is created during setup, or you can create one manually:

    ```bash
    cp .env.example .env  # If an example exists, otherwise create .env
    ```

    Ensure your `.env` contains:
    ```env
    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=transfers
    DB_SSL_MODE=disable
    SERVER_PORT=8080
    SERVER_HOST=0.0.0.0
    LOG_LEVEL=info
    ```

## Running the Application

### Using Docker (Recommended)

To start the application and all dependencies (PostgreSQL):

```bash
make docker-up
```

This command will:
*   Start the PostgreSQL database and the API server in Docker containers.
*   Wait for the database to be ready.
*   Run database migrations automatically.

To stop the services:

```bash
make docker-down
```

### Running Locally

1.  **Start the Database:**

    You need a running PostgreSQL instance. You can use Docker for just the DB:
    ```bash
    docker-compose up -d postgres
    ```

2.  **Run Migrations:**

    ```bash
    make migrate-up
    ```

3.  **Run the Server:**

    ```bash
    make run
    ```

    The server will start at `http://localhost:8080`.

## API Endpoints

*   **POST /accounts**: Create a new account
    ```json
    {
      "account_id": 1,
      "initial_balance": "100.50"
    }
    ```

*   **GET /accounts/{account_id}**: Get account details

*   **POST /transactions**: Transfer funds
    ```json
    {
      "source_account_id": 1,
      "destination_account_id": 2,
      "amount": "50.00"
    }
    ```

