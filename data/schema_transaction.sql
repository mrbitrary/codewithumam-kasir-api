CREATE TABLE IF NOT EXISTS core.transaction (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    total_items INT NOT NULL,
    total_price_amount BIGINT NOT NULL,
    total_price_scale INT NOT NULL DEFAULT 2,
    total_price_display NUMERIC(18, 8) GENERATED ALWAYS AS (
        total_price_amount::numeric / (10 ^ total_price_scale)::numeric
    ) STORED,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS core.transaction_detail (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    transaction_id UUID NOT NULL REFERENCES core.transaction(id) ON DELETE CASCADE,
    product_id UUID REFERENCES core.product(id) ON DELETE SET NULL,
    product_name TEXT NOT NULL,
    category_id UUID REFERENCES core.category(id) ON DELETE SET NULL,
    category_name TEXT NOT NULL,
    price_amount BIGINT NOT NULL,
    price_scale INT NOT NULL,
    price_display NUMERIC(18, 8) GENERATED ALWAYS AS (
        price_amount::numeric / (10 ^ price_scale)::numeric
    ) STORED,
    currency VARCHAR(3) NOT NULL,
    quantity INT NOT NULL,
    total_price_amount BIGINT NOT NULL,
    total_price_scale INT NOT NULL,
    total_price_display NUMERIC(18, 8) GENERATED ALWAYS AS (
        total_price_amount::numeric / (10 ^ total_price_scale)::numeric
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_transaction_date ON core.transaction(created_at);
CREATE INDEX idx_transaction_detail_product ON core.transaction_detail(product_id);
CREATE INDEX idx_transaction_detail_category ON core.transaction_detail(category_id);

CREATE TRIGGER trg_transaction_version_increment
BEFORE UPDATE ON core.transaction
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();

CREATE TRIGGER trg_transaction_detail_version_increment
BEFORE UPDATE ON core.transaction_detail
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();
