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
    record_uuid UUID;
BEGIN
    -- Handle special tenant derivation for tables without direct tenant_id
    IF TG_TABLE_NAME = 'tenant_user_roles' THEN
        -- Get tenant_id from the related tenant_user record
        SELECT tu.tenant_id INTO tenant 
        FROM tenant_users tu 
        WHERE tu.id = COALESCE(NEW.tenant_user_id, OLD.tenant_user_id);
        record_uuid := NULL; -- Composite key table, no single UUID
    ELSE
        tenant := current_tenant_id();
        record_uuid := COALESCE(NEW.id, OLD.id);
    END IF;

    IF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, new_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, record_uuid, 'INSERT', row_to_json(NEW)::jsonb);
        RETURN NEW;

    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, old_data, new_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, record_uuid, 'UPDATE', row_to_json(OLD)::jsonb, row_to_json(NEW)::jsonb);
        RETURN NEW;

    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_logs(tenant_id, user_id, table_name, record_id, action, old_data)
        VALUES (tenant, current_setting('app.current_user', true)::UUID, TG_TABLE_NAME, record_uuid, 'DELETE', row_to_json(OLD)::jsonb);
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
        AND tablename IN ('tenant_users','teachers','students','subjects','classes','class_subjects','academic_years','attendance', 'grades','enrollments','schedules','departments','parents','notifications','student_fees','fee_types','tenant_features','tenant_user_roles')
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
