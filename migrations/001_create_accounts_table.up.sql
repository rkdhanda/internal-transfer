CREATE TABLE IF NOT EXISTS accounts (
    account_id BIGINT PRIMARY KEY,
    balance DECIMAL(36, 18) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_balance ON accounts(balance);
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at DESC);
