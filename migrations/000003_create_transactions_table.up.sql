CREATE TABLE IF NOT EXISTS transactions (
                                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                            account_id UUID NOT NULL CONSTRAINT fk_account REFERENCES accounts ON DELETE CASCADE,
                                            transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
                                            amount NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
                                            balance_before NUMERIC(18, 2) NOT NULL,
                                            balance_after NUMERIC(18, 2) NOT NULL,
                                            description TEXT,
                                            transaction_date TIMESTAMP NOT NULL,
                                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);