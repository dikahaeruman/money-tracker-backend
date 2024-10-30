CREATE TABLE Accounts (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          user_id integer NOT NULL,
                          account_name VARCHAR(255) NOT NULL,
                          balance DECIMAL(18, 2) DEFAULT 0.00,
                          currency VARCHAR(10) DEFAULT 'IDR',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
);
