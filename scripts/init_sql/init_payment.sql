CREATE TABLE IF NOT EXISTS balances
(
    user_id  TEXT PRIMARY KEY,
    balance  BIGINT NOT NULL DEFAULT 0,
    currency TEXT   NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions
(
    id               SERIAL PRIMARY KEY,
    booking_id       INTEGER,
    user_id          TEXT    NOT NULL,
    amount           BIGINT  NOT NULL,
    transaction_type TEXT    NOT NULL CHECK ( transaction_type IN ('charge', 'payment', 'refund', 'promocode_activate', 'promocode_generate') ),
    status           TEXT    NOT NULL CHECK ( status IN ('pending', 'completed', 'failed', 'canceled') ) DEFAULT 'pending',
    description      TEXT,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS promocodes
(
    code           TEXT PRIMARY KEY,
    amount         BIGINT  NOT NULL,
    max_uses       INTEGER NOT NULL DEFAULT 1,
    used_count     INTEGER NOT NULL DEFAULT 0,
    expires_at     TIMESTAMP,
    created_by     TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_promocodes_expires_at ON promocodes(expires_at);
CREATE INDEX IF NOT EXISTS idx_promocodes_created_by ON promocodes(created_by);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_booking_id ON transactions(booking_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

CREATE OR REPLACE FUNCTION update_balance_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_balance_updated_at
    BEFORE UPDATE ON balances
    FOR EACH ROW
    EXECUTE FUNCTION update_balance_updated_at();

CREATE TRIGGER trigger_transaction_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_balance_updated_at();

