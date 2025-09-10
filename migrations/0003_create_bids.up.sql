BEGIN;

-- bids
CREATE TABLE IF NOT EXISTS bids (
                                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    executor_id UUID NOT NULL REFERENCES users(id),
    cover_text TEXT,
    price BIGINT NULL, -- proposed price (in currency minor units)
    proposed_deadline TIMESTAMP WITH TIME ZONE NULL,
                                                     attachments JSONB DEFAULT '[]'::jsonb,
                                                     status VARCHAR(32) NOT NULL DEFAULT 'created', -- created|pending_payment|paid|visible_to_client|shortlisted|won|lost
    paid_at TIMESTAMP WITH TIME ZONE,
    visibility_to_client BOOLEAN DEFAULT FALSE,
    metadata JSONB DEFAULT '{}'::jsonb, -- for future fields like contact_purchase_id
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_bids_order_status ON bids (order_id, status);
CREATE INDEX IF NOT EXISTS idx_bids_executor ON bids (executor_id);
CREATE INDEX IF NOT EXISTS idx_bids_paid_at ON bids (paid_at);

COMMIT;
