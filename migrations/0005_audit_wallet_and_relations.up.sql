BEGIN;

-- wallet transactions (ledger-style)
CREATE TABLE IF NOT EXISTS wallet_transactions (
                                                   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    organization_id UUID NULL REFERENCES organizations(id),
    payment_id UUID NULL REFERENCES payments(id),
    amount BIGINT NOT NULL, -- positive for credit, negative for debit
    balance_snapshot BIGINT NULL, -- optional snapshot after operation (if computing)
    type VARCHAR(64) NOT NULL, -- credit|debit|refund|fee
    meta JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_wallet_user_created ON wallet_transactions (user_id, created_at DESC);

-- simple audit log
CREATE TABLE IF NOT EXISTS audit_logs (
                                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id UUID NULL REFERENCES users(id),
    action VARCHAR(128) NOT NULL,
    object_type VARCHAR(64),
    object_id UUID,
    payload JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_audit_actor ON audit_logs (actor_id);

COMMIT;
