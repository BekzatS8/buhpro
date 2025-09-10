BEGIN;

-- orders
CREATE TABLE IF NOT EXISTS orders (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    client_user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(512) NOT NULL,
    description TEXT,
    category VARCHAR(128),
    subcategory VARCHAR(128),
    region VARCHAR(128),
    mode_online BOOLEAN DEFAULT TRUE, -- online/offline
    deadline TIMESTAMP WITH TIME ZONE,
                                                          budget_min BIGINT, -- in KZT (or currency minor units)
                                                          budget_max BIGINT,
                                                          currency VARCHAR(8) DEFAULT 'KZT',
    status VARCHAR(32) NOT NULL DEFAULT 'draft', -- draft|pending_payment|published|executor_selected|in_progress|client_review|completed|archived|cancelled
    promotion_flags JSONB DEFAULT '{}'::jsonb, -- {top_until: ts, highlighted: bool, pinned_until: ts, base_price: int}
    attachments JSONB DEFAULT '[]'::jsonb, -- list of file meta
    chosen_bid_id UUID NULL, -- FK not enforced to avoid circular constraints; set by app
    meta JSONB DEFAULT '{}'::jsonb, -- future-proof
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    published_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

-- indexes for feed and filters
CREATE INDEX IF NOT EXISTS idx_orders_status_category_region_created ON orders (status, category, region, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_deadline ON orders (deadline);
CREATE INDEX IF NOT EXISTS idx_orders_published_at ON orders (published_at);

COMMIT;
