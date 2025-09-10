BEGIN;

-- payments
CREATE TABLE IF NOT EXISTS payments (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NULL REFERENCES users(id), -- payer
    organization_id UUID NULL REFERENCES organizations(id), -- if payment from org
    related_type VARCHAR(32) NOT NULL, -- 'order_publish' | 'bid_fee' | 'buy_contact' | 'wallet_topup' | ...
    related_id UUID NULL, -- order_id or bid_id etc.
    provider VARCHAR(64) NOT NULL, -- kaspi|halyk|stripe|manual
    provider_payment_id VARCHAR(255) NULL, -- id from provider
    amount BIGINT NOT NULL, -- in minor units
    currency VARCHAR(8) DEFAULT 'KZT',
    status VARCHAR(32) NOT NULL DEFAULT 'initiated', -- initiated|redirected|success|failed|expired|refunded
    items JSONB DEFAULT '[]'::jsonb, -- details of items paid
    idempotency_key VARCHAR(255), -- stored for safe retries
    expires_at TIMESTAMP WITH TIME ZONE,
    webhook_meta JSONB DEFAULT '{}'::jsonb, -- raw provider response
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_provider_pid ON payments (provider, provider_payment_id) WHERE provider_payment_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_user ON payments (user_id);
CREATE INDEX IF NOT EXISTS idx_payments_idempotency ON payments (idempotency_key);

COMMIT;
