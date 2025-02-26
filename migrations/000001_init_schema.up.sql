-- Создание таблицы кошельков
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    address VARCHAR(100) UNIQUE NOT NULL,
    balance NUMERIC(20,2) DEFAULT 100.00 CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы транзакций
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_address VARCHAR(100) NOT NULL REFERENCES wallets(address),
    receiver_address VARCHAR(100) NOT NULL REFERENCES wallets(address),
    amount NUMERIC(20,2) NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации запросов
CREATE INDEX idx_wallets_address ON wallets(address);
CREATE INDEX idx_transactions_sender ON transactions(sender_address);
CREATE INDEX idx_transactions_receiver ON transactions(receiver_address);
CREATE INDEX idx_transactions_timestamp ON transactions(timestamp);