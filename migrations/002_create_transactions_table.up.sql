CREATE TABLE IF NOT EXISTS transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    source_account_id BIGINT NOT NULL,
    destination_account_id BIGINT NOT NULL,
    amount DECIMAL(36, 18) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    error_message TEXT,
    
    CONSTRAINT fk_source_account 
        FOREIGN KEY (source_account_id) 
        REFERENCES accounts(account_id),
    CONSTRAINT fk_destination_account 
        FOREIGN KEY (destination_account_id) 
        REFERENCES accounts(account_id),
    CONSTRAINT check_different_accounts 
        CHECK (source_account_id != destination_account_id)
);

CREATE INDEX IF NOT EXISTS idx_transactions_source ON transactions(source_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_destination ON transactions(destination_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
