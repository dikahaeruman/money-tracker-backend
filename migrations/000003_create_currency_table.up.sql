-- Create Currency table
CREATE TABLE IF NOT EXISTS Currencies (
    id SERIAL PRIMARY KEY,
    currency_code VARCHAR(10) NOT NULL UNIQUE,
    currency_name VARCHAR(255) NOT NULL
);

-- Alter Accounts table to reference Currencies
ALTER TABLE Accounts
DROP COLUMN currency,
ADD COLUMN currency_id INTEGER REFERENCES Currencies(id);
