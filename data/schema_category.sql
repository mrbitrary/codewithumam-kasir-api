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
CREATE TRIGGER trg_category_version_increment
BEFORE UPDATE ON core.category
FOR EACH ROW EXECUTE FUNCTION core.fn_increment_version();
