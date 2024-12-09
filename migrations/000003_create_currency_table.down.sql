-- Rollback: Add back currency column and drop currency_id column
ALTER TABLE Accounts
DROP COLUMN currency_id,
ADD COLUMN currency VARCHAR(10) DEFAULT 'IDR';

-- Rollback: Drop Currencies table
DROP TABLE IF EXISTS Currencies;
