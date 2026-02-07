CREATE TABLE IF NOT EXISTS core.sales_summary_daily (
    report_date DATE NOT NULL,
    product_id UUID NOT NULL,
    category_id UUID NOT NULL, -- Use 00000000-0000-0000-0000-000000000000 for Uncategorized
    total_sold INT NOT NULL DEFAULT 0,
    total_revenue BIGINT NOT NULL DEFAULT 0,
    
    PRIMARY KEY (report_date, product_id, category_id)
);

CREATE TABLE IF NOT EXISTS core.transaction_summary_daily (
    report_date DATE PRIMARY KEY,
    total_revenue BIGINT NOT NULL DEFAULT 0,
    total_transactions INT NOT NULL DEFAULT 0
);

CREATE OR REPLACE FUNCTION core.fn_update_sales_summary()
RETURNS TRIGGER AS $$
DECLARE
    v_category_id UUID;
BEGIN
    -- Handle category_id if NULL for NEW record
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        v_category_id := COALESCE(NEW.category_id, '00000000-0000-0000-0000-000000000000'::uuid);
    END IF;

    -- Handle INSERT
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO core.sales_summary_daily (report_date, product_id, category_id, total_sold, total_revenue)
        VALUES (
            DATE(NEW.created_at), 
            NEW.product_id, 
            v_category_id, 
            NEW.quantity, 
            NEW.total_price_amount
        )
        ON CONFLICT (report_date, product_id, category_id) DO UPDATE SET
            total_sold = core.sales_summary_daily.total_sold + EXCLUDED.total_sold,
            total_revenue = core.sales_summary_daily.total_revenue + EXCLUDED.total_revenue;
        RETURN NEW;
    
    -- Handle DELETE
    ELSIF (TG_OP = 'DELETE') THEN
        v_category_id := COALESCE(OLD.category_id, '00000000-0000-0000-0000-000000000000'::uuid);
        UPDATE core.sales_summary_daily
        SET 
            total_sold = total_sold - OLD.quantity,
            total_revenue = total_revenue - OLD.total_price_amount
        WHERE 
            report_date = DATE(OLD.created_at) 
            AND product_id = OLD.product_id 
            AND category_id = v_category_id;
        RETURN OLD;

    -- Handle UPDATE
    ELSIF (TG_OP = 'UPDATE') THEN
        -- Decrement old values
        DECLARE
            v_old_category_id UUID;
        BEGIN
            v_old_category_id := COALESCE(OLD.category_id, '00000000-0000-0000-0000-000000000000'::uuid);
            UPDATE core.sales_summary_daily
            SET 
                total_sold = total_sold - OLD.quantity,
                total_revenue = total_revenue - OLD.total_price_amount
            WHERE 
                report_date = DATE(OLD.created_at) 
                AND product_id = OLD.product_id 
                AND category_id = v_old_category_id;
        END;
        
        -- Increment new values
        INSERT INTO core.sales_summary_daily (report_date, product_id, category_id, total_sold, total_revenue)
        VALUES (
            DATE(NEW.created_at), 
            NEW.product_id, 
            v_category_id, 
            NEW.quantity, 
            NEW.total_price_amount
        )
        ON CONFLICT (report_date, product_id, category_id) DO UPDATE SET
            total_sold = core.sales_summary_daily.total_sold + EXCLUDED.total_sold,
            total_revenue = core.sales_summary_daily.total_revenue + EXCLUDED.total_revenue;
        RETURN NEW;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION core.fn_update_transaction_summary()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO core.transaction_summary_daily (report_date, total_revenue, total_transactions)
        VALUES (DATE(NEW.created_at), NEW.total_price_amount, 1)
        ON CONFLICT (report_date) DO UPDATE SET
            total_revenue = core.transaction_summary_daily.total_revenue + EXCLUDED.total_revenue,
            total_transactions = core.transaction_summary_daily.total_transactions + EXCLUDED.total_transactions;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE core.transaction_summary_daily
        SET 
            total_revenue = total_revenue - OLD.total_price_amount,
            total_transactions = total_transactions - 1
        WHERE report_date = DATE(OLD.created_at);
    ELSIF (TG_OP = 'UPDATE') THEN
        UPDATE core.transaction_summary_daily
        SET total_revenue = total_revenue - OLD.total_price_amount + NEW.total_price_amount
        WHERE report_date = DATE(NEW.created_at);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_refresh_sales_summary
AFTER INSERT OR UPDATE OR DELETE ON core.transaction_detail
FOR EACH ROW EXECUTE FUNCTION core.fn_update_sales_summary();

CREATE TRIGGER trg_refresh_transaction_summary
AFTER INSERT OR UPDATE OR DELETE ON core.transaction
FOR EACH ROW EXECUTE FUNCTION core.fn_update_transaction_summary();
