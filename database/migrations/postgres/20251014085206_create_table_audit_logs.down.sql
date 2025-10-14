-- =========================================
-- ROLLBACK AUDIT LOGS MIGRATION
-- =========================================

-- =========================================
-- DROP AUDIT TRIGGERS
-- =========================================
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT tablename FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename IN ('users','teachers','students','subjects','classes','class_subjects','academic_years','attendance', 'grades','enrollments','schedules','departments','parents','notifications','student_fees','fee_types','tenant_features')
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS trg_audit_%1$I ON %1$I;', tbl);
    END LOOP;
END;
$$;

-- =========================================
-- DROP AUDIT FUNCTION
-- =========================================
DROP FUNCTION IF EXISTS fn_audit_log();

-- =========================================
-- DROP RLS POLICY AND DISABLE RLS
-- =========================================
DROP POLICY IF EXISTS tenant_isolation ON audit_logs;
ALTER TABLE IF EXISTS audit_logs DISABLE ROW LEVEL SECURITY;

-- =========================================
-- DROP AUDIT LOGS TABLE
-- =========================================
DROP TABLE IF EXISTS audit_logs;
