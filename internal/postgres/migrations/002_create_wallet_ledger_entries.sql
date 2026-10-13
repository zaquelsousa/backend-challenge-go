-- +goose Up
CREATE TABLE wallet_ledger_entries (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL
        REFERENCES wallets(id) ON DELETE RESTRICT,

    transaction_id UUID NOT NULL,
    direction TEXT NOT NULL
        CHECK (direction IN ('CREDIT', 'DEBIT')),

    amount_minor BIGINT NOT NULL
        CHECK (amount_minor > 0),

    currency CHAR(3) NOT NULL,

    balance_before_minor BIGINT NOT NULL
        CHECK (balance_before_minor >= 0),

    balance_after_minor BIGINT NOT NULL
        CHECK (balance_after_minor >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT wallet_ledger_currency_not_empty
        CHECK (currency = UPPER(currency)),

    CONSTRAINT wallet_ledger_balance_transition
        CHECK (
            (direction = 'CREDIT'
                AND balance_after_minor =
                    balance_before_minor + amount_minor)
            OR
            (direction = 'DEBIT'
                AND balance_before_minor >= amount_minor
                AND balance_after_minor =
                    balance_before_minor - amount_minor)
        ),

    CONSTRAINT wallet_ledger_wallet_transaction_unique
        UNIQUE (wallet_id, transaction_id)
);

CREATE INDEX idx_wallet_ledger_wallet_created
    ON wallet_ledger_entries (wallet_id, created_at DESC, id DESC);


-- +goose Down
DROP TABLE wallet_ledger_entries;
