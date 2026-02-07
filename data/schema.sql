-- CREATE DATABASE kasir_api;
---
CREATE SCHEMA IF NOT EXISTS core;
CREATE SCHEMA IF NOT EXISTS audit;
--- core schema
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
