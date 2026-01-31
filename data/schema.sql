CREATE DATABASE kasir_api;
---
CREATE SCHEMA IF NOT EXISTS core;
CREATE SCHEMA IF NOT EXISTS audit;
---

--- core schema
CREATE TABLE IF NOT EXISTS core.category (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    version    INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,

    name TEXT NOT NULL,
    description TEXT,

    CONSTRAINT name_not_empty CHECK (char_length(trim(name)) > 0)
);
---
CREATE INDEX idx_category_deleted ON core.category (deleted_at)
WHERE deleted_at IS NOT NULL;
---
CREATE UNIQUE INDEX idx_category_active_name ON core.category (lower(name))
WHERE deleted_at IS NULL;
---
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

    category_id UUID REFERENCES core.category(id) ON DELETE SET NULL
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
CREATE OR REPLACE FUNCTION core.fn_increment_version()
RETURNS TRIGGER AS $$
DECLARE
    v_user_id TEXT;
BEGIN
    v_user_id := coalesce(current_setting('app.current_user_id', true), 'SYSTEM');

    IF (NEW IS DISTINCT FROM OLD) OR (OLD.deleted_at IS DISTINCT FROM NEW.deleted_at) THEN
        NEW.version = OLD.version + 1;
        NEW.updated_at = CURRENT_TIMESTAMP;
        NEW.updated_by = v_user_id;
    END IF;
    RETURN NEW; -- BEFORE trigger must return NEW/OLD/NULL to control operation
END;
$$ LANGUAGE plpgsql;
---
CREATE TRIGGER trg_product_version_increment
BEFORE UPDATE ON core.product
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();
---
CREATE TRIGGER trg_category_version_increment
BEFORE UPDATE ON core.category
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();
---

--- audit schema
CREATE TABLE audit.audit_log (
    id            UUID PRIMARY KEY DEFAULT uuidv7(),
    table_name    TEXT NOT NULL,
    operation     TEXT NOT NULL, -- INSERT, UPDATE, DELETE
    row_id        TEXT NOT NULL, -- The ID of the record being changed
    changed_by    TEXT,          -- Captured from a session variable
    old_data      JSONB,         -- The row BEFORE the change (NULL for INSERT)
    new_data      JSONB,         -- The row AFTER the change (NULL for DELETE)
    changed_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_operation CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE'))
);
---
CREATE INDEX idx_audit_log_table_operation ON audit.audit_log (table_name, operation);
CREATE INDEX idx_audit_log_table_row ON audit.audit_log (table_name, row_id);
CREATE INDEX idx_audit_log_table_changed_at ON audit.audit_log (table_name, changed_at DESC);
CREATE INDEX idx_audit_log_table_changed_by ON audit.audit_log (table_name, changed_by);
---
CREATE OR REPLACE FUNCTION audit.fn_audit_log()
RETURNS TRIGGER AS $$
DECLARE
    v_old_data JSONB := NULL;
    v_new_data JSONB := NULL;
    v_user_id TEXT;
BEGIN
    v_user_id := coalesce(current_setting('app.current_user_id', true), 'SYSTEM');

    IF (TG_OP = 'UPDATE') THEN
        v_old_data := to_jsonb(OLD);
        v_new_data := to_jsonb(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        v_old_data := to_jsonb(OLD);
    ELSIF (TG_OP = 'INSERT') THEN
        v_new_data := to_jsonb(NEW);
    END IF;

    INSERT INTO audit.audit_log (table_name, operation, row_id, changed_by, old_data, new_data)
    VALUES (
        TG_TABLE_NAME,
        TG_OP,
        COALESCE(NEW.id::TEXT, OLD.id::TEXT),
        v_user_id,
        v_old_data,
        v_new_data
    );

    RETURN NULL; -- For AFTER triggers, return value doesn't matter
END;
$$ LANGUAGE plpgsql;
---
CREATE TRIGGER trg_audit_product
AFTER INSERT OR UPDATE OR DELETE ON core.product
FOR EACH ROW EXECUTE FUNCTION audit.fn_audit_log();
---
CREATE TRIGGER trg_audit_category
AFTER INSERT OR UPDATE OR DELETE ON core.category
FOR EACH ROW EXECUTE FUNCTION audit.fn_audit_log();
---
