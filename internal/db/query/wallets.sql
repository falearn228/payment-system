-- name: GetWalletForUpdate :one
SELECT balance FROM wallets 
WHERE address = $1 
FOR UPDATE;

-- name: DecrementBalance :exec
UPDATE wallets 
SET balance = balance - $1,
    updated_at = CURRENT_TIMESTAMP 
WHERE address = $2;

-- name: IncrementBalance :exec
UPDATE wallets 
SET balance = balance + $1,
    updated_at = CURRENT_TIMESTAMP 
WHERE address = $2;

-- name: CreateTransaction :exec
INSERT INTO transactions (sender_address, receiver_address, amount) 
VALUES ($1, $2, $3);

-- name: GetBalance :one
SELECT balance FROM wallets 
WHERE address = $1;

-- name: CreateWallet :exec
INSERT INTO wallets (address) 
VALUES (gen_random_uuid()::VARCHAR(100));