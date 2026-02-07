CREATE TABLE IF NOT EXISTS core.product (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    version    INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,

    name TEXT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    price_amount BIGINT NOT NULL,
    price_scale INT NOT NULL,
    currency VARCHAR(3) NOT NULL,

    price_display NUMERIC(18, 8) GENERATED ALWAYS AS (
        price_amount::numeric / (10 ^ price_scale)::numeric
    ) STORED,

    CONSTRAINT name_not_empty CHECK (char_length(trim(name)) > 0),
    CONSTRAINT price_not_negative CHECK (price_amount >= 0),
    CONSTRAINT currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT scale_range CHECK (price_scale >= 0 AND price_scale <= 8),
    CONSTRAINT stock_not_negative CHECK (stock >= 0),

    category_id UUID REFERENCES core.category(id) ON DELETE SET NULL,

    name_tsvector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', name)
    ) STORED
);
---
CREATE INDEX idx_product_deleted ON core.product (deleted_at)
WHERE deleted_at IS NOT NULL;
---
CREATE INDEX idx_product_active_name ON core.product (lower(name))
WHERE deleted_at IS NULL;
---
CREATE INDEX idx_product_active_price ON core.product (price_display)
WHERE deleted_at IS NULL;
---
CREATE INDEX idx_product_category_id ON core.product (category_id)
WHERE deleted_at IS NULL;
---
CREATE INDEX idx_product_category_price ON core.product (category_id, price_display)
WHERE deleted_at IS NULL;
---
CREATE INDEX idx_product_name_tsvector ON core.product USING GIN (name_tsvector);
---
CREATE TRIGGER trg_product_version_increment
BEFORE UPDATE ON core.product
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();
