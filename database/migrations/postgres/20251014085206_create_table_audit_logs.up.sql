-- =========================================
-- AUDIT LOGS
-- =========================================

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    table_name VARCHAR(255) NOT NULL,
    record_id UUID,
    action VARCHAR(50) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    old_data JSONB,
    new_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Enable Row-Level Security
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON audit_logs
    USING (tenant_id = current_tenant_id());

-- =========================================
-- TRIGGER FUNCTION FOR AUDIT LOGGING
-- =========================================

CREATE OR REPLACE FUNCTION fn_audit_log() RETURNS TRIGGER AS $$
DECLARE
    tenant UUID;
BEGIN
    tenant := current_tenant_id();

    IF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, new_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, NEW.id, 'INSERT', row_to_json(NEW)::jsonb);
        RETURN NEW;

    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, old_data, new_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, NEW.id, 'UPDATE', row_to_json(OLD)::jsonb, row_to_json(NEW)::jsonb);
        RETURN NEW;

    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, old_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, OLD.id, 'DELETE', row_to_json(OLD)::jsonb);
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- APPLY AUDIT TRIGGERS
-- =========================================
-- You can choose which tables to log.

DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT tablename FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename IN ('users','teachers','students','subjects','classes','class_subjects','academic_years','attendance', 'grades','enrollments','schedules','departments','parents','notifications','student_fees','fee_types','tenant_features')
    LOOP
        EXECUTE format(
            'CREATE TRIGGER trg_audit_%1$I
             AFTER INSERT OR UPDATE OR DELETE ON %1$I
             FOR EACH ROW EXECUTE FUNCTION fn_audit_log();',
            tbl
        );
    END LOOP;
END;
$$;

-- =========================================
-- ADDITIONAL SETTINGS
-- =========================================
-- Your app should set both tenant and user context before each transaction:
--   SELECT set_config('app.current_tenant', '<tenant_uuid>', false);
--   SELECT set_config('app.current_user', '<user_uuid>', false);
