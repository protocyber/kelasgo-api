-- ======================================================
-- SCHOOL MANAGEMENT SYSTEM DATABASE SCHEMA
-- PostgreSQL version
-- Author: Fitrah
-- ======================================================
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ======================================================
-- ENUMS
-- ======================================================
CREATE TYPE fee_status_enum AS ENUM ('paid', 'unpaid', 'partial', 'overdue');
CREATE TYPE subscription_status_enum AS ENUM ('active', 'inactive', 'cancelled', 'expired', 'trial');
CREATE TYPE subscription_plan_status_enum AS ENUM ('active', 'inactive', 'cancelled', 'expired');
CREATE TYPE invoice_status_enum AS ENUM ('draft', 'sent', 'paid', 'unpaid', 'overdue', 'cancelled');
CREATE TYPE attendance_status_enum AS ENUM ('Present', 'Absent', 'Late', 'Excused');

-- ======================================================
-- SUBSCRIPTION PLANS
-- ======================================================
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(100) NOT NULL,
    price_monthly DECIMAL(10,2) DEFAULT 0,
    price_yearly DECIMAL(10,2) DEFAULT 0,
    max_students INT,
    max_teachers INT,
    storage_limit_mb INT,
    features JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =======================
-- TENANTS
-- =======================
CREATE TABLE
  tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    contact_email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    logo_url VARCHAR(255),
    plan_id UUID,
    subscription_status subscription_status_enum DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

-- ======================================================
-- SUBSCRIPTIONS (per tenant)
-- ======================================================
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID,
    plan_id UUID,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_trial BOOLEAN DEFAULT FALSE,
    status subscription_plan_status_enum DEFAULT 'active',
    amount_paid DECIMAL(10,2),
    payment_method VARCHAR(50),
    invoice_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ======================================================
-- INVOICES (billing)
-- ======================================================
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID,
    subscription_id UUID,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'Rp',
    issue_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    status invoice_status_enum DEFAULT 'unpaid',
    payment_date DATE,
    payment_reference VARCHAR(100)
);

-- =======================
-- FEATURE FLAGS
-- =======================
CREATE TABLE
  feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT
  );

INSERT INTO
  feature_flags (code, name, description)
VALUES
  (
    'academic_years',
    'Academic Years',
    'Allow schools to manage multiple academic years'
  ),
  (
    'attendance_tracking',
    'Attendance Tracking',
    'Track daily student attendance'
  ),
  (
    'grading_system',
    'Grading System',
    'Enable advanced grading system per subject'
  );

-- =======================
-- TENANT FEATURES
-- =======================
CREATE TABLE
  tenant_features (
    tenant_id UUID,
    feature_id UUID,
    enabled BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (tenant_id, feature_id)
  );

-- ======================================================
-- ROLES
-- ======================================================
CREATE TABLE
  roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
  );

-- ======================================================
-- USERS
-- ======================================================
CREATE TABLE
  users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('Male', 'Female')),
    date_of_birth DATE,
    phone VARCHAR(20),
    address TEXT,
    is_active BOOLEAN DEFAULT TRUE
  );

-- ======================================================
-- TENANT USERS (many-to-many relationship between users and tenants)
-- ======================================================
CREATE TABLE
  tenant_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, user_id)
  );

-- =======================
-- USER ROLES
-- =======================
CREATE TABLE
  user_roles (
    user_id UUID,
    role_id UUID,
    PRIMARY KEY (user_id, role_id)
  );

-- ======================================================
-- DEPARTMENTS
-- ======================================================
CREATE TABLE
  departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    head_teacher_id UUID
  );

-- ======================================================
-- TEACHERS
-- ======================================================
CREATE TABLE
  teachers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    tenant_user_id UUID NOT NULL,
    employee_number VARCHAR(50) UNIQUE,
    hire_date DATE,
    department_id UUID,
    qualification VARCHAR(100),
    position VARCHAR(100)
  );

-- ======================================================
-- PARENTS
-- ======================================================
CREATE TABLE
  parents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    relationship VARCHAR(50)
  );

-- ======================================================
-- SCHOOL YEARS
-- ======================================================
CREATE TABLE
  academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE
  );

-- ======================================================
-- CLASSES
-- ======================================================
CREATE TABLE
  classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    grade_level INT,
    homeroom_teacher_id UUID,
    academic_year_id UUID
  );

-- ======================================================
-- STUDENTS
-- ======================================================
CREATE TABLE
  students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    tenant_user_id UUID NOT NULL,
    student_number VARCHAR(50) UNIQUE NOT NULL,
    admission_date DATE NOT NULL,
    class_id UUID,
    parent_id UUID
  );

-- ======================================================
-- SUBJECTS
-- ======================================================
CREATE TABLE
  subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    department_id UUID,
    credit INT DEFAULT 0
  );

-- ======================================================
-- CLASS SUBJECTS (link between class, subject, teacher)
-- ======================================================
CREATE TABLE
  class_subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    class_id UUID,
    subject_id UUID,
    teacher_id UUID,
    UNIQUE (class_id, subject_id, teacher_id)
  );

-- ======================================================
-- SCHEDULES
-- ======================================================
CREATE TABLE
  schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    class_subject_id UUID,
    day_of_week VARCHAR(15) CHECK (
      day_of_week IN (
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday'
      )
    ),
    start_time TIME,
    end_time TIME,
    room VARCHAR(50)
  );

-- ======================================================
-- ENROLLMENTS
-- ======================================================
CREATE TABLE
  enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    student_id UUID,
    class_subject_id UUID,
    academic_year_id UUID
  );

-- ======================================================
-- GRADES
-- ======================================================
CREATE TABLE
  grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    enrollment_id UUID,
    grade_type VARCHAR(50) CHECK (
      grade_type IN ('Assignment', 'Midterm', 'Final', 'Other')
    ),
    score DECIMAL(5, 2),
    remarks TEXT
  );

-- ======================================================
-- ATTENDANCE
-- ======================================================
CREATE TABLE
  attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    student_id UUID,
    schedule_id UUID,
    status attendance_status_enum,
    attendance_date DATE DEFAULT CURRENT_DATE,
    remarks TEXT
  );

-- ======================================================
-- NOTIFICATIONS
-- ======================================================
CREATE TABLE
  notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    user_id UUID,
    title VARCHAR(100),
    message TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

-- ========================================
-- fee_types
-- ========================================
CREATE TABLE
  fee_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    default_amount DECIMAL(10, 2) DEFAULT 0,
    is_mandatory BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT chk_fee_types_amount CHECK (default_amount >= 0)
  );

-- ========================================
-- student_fees
-- ========================================
CREATE TABLE
  student_fees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    tenant_id UUID NOT NULL,
    student_id UUID,
    fee_type_id UUID,
    academic_year_id UUID,
    amount DECIMAL(10, 2) NOT NULL,
    due_date DATE NOT NULL,
    status fee_status_enum DEFAULT 'unpaid',
    payment_date DATE,
    payment_method VARCHAR(50),
    notes TEXT,
    CONSTRAINT chk_student_fees_amount CHECK (amount >= 0),
    CONSTRAINT chk_student_fees_payment_date CHECK (
      payment_date IS NULL
      OR payment_date >= due_date
      OR status IN ('paid', 'partial')
    )
  );

-- ======================================================
-- FOREIGN KEY CONSTRAINTS FOR BUSINESS LOGIC
-- ======================================================
-- Add foreign key constraints for business relationships
-- Moved outside table definitions to avoid circular dependencies

-- Subscription and billing foreign key constraints
ALTER TABLE tenants ADD CONSTRAINT fk_tenants_plan_id FOREIGN KEY (plan_id) REFERENCES subscription_plans (id) ON DELETE SET NULL;

ALTER TABLE subscriptions ADD CONSTRAINT fk_subscriptions_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_subscriptions_plan_id FOREIGN KEY (plan_id) REFERENCES subscription_plans (id) ON DELETE CASCADE;

ALTER TABLE invoices ADD CONSTRAINT fk_invoices_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_invoices_subscription_id FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE;

-- Tenant foreign key constraints
ALTER TABLE tenant_features ADD CONSTRAINT fk_tenant_features_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_tenant_features_feature_id FOREIGN KEY (feature_id) REFERENCES feature_flags (id) ON DELETE CASCADE;

ALTER TABLE tenant_users ADD CONSTRAINT fk_tenant_users_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_tenant_users_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE;

ALTER TABLE departments ADD CONSTRAINT fk_departments_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE;

ALTER TABLE teachers ADD CONSTRAINT fk_teachers_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_teachers_tenant_user_id FOREIGN KEY (tenant_user_id) REFERENCES tenant_users (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_teachers_department_id FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE SET NULL;

ALTER TABLE departments ADD CONSTRAINT fk_departments_head_teacher_id FOREIGN KEY (head_teacher_id) REFERENCES teachers (id) ON DELETE SET NULL;

ALTER TABLE parents ADD CONSTRAINT fk_parents_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE;

ALTER TABLE academic_years ADD CONSTRAINT fk_academic_years_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE;

ALTER TABLE classes ADD CONSTRAINT fk_classes_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_classes_homeroom_teacher_id FOREIGN KEY (homeroom_teacher_id) REFERENCES teachers (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_classes_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE SET NULL;

ALTER TABLE students ADD CONSTRAINT fk_students_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_students_tenant_user_id FOREIGN KEY (tenant_user_id) REFERENCES tenant_users (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_students_class_id FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_students_parent_id FOREIGN KEY (parent_id) REFERENCES parents (id) ON DELETE SET NULL;

ALTER TABLE subjects ADD CONSTRAINT fk_subjects_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_subjects_department_id FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE SET NULL;

ALTER TABLE class_subjects ADD CONSTRAINT fk_class_subjects_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_class_subjects_class_id FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_class_subjects_subject_id FOREIGN KEY (subject_id) REFERENCES subjects (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_class_subjects_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers (id) ON DELETE SET NULL;

ALTER TABLE schedules ADD CONSTRAINT fk_schedules_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_schedules_class_subject_id FOREIGN KEY (class_subject_id) REFERENCES class_subjects (id) ON DELETE CASCADE;

ALTER TABLE enrollments ADD CONSTRAINT fk_enrollments_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_enrollments_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_enrollments_class_subject_id FOREIGN KEY (class_subject_id) REFERENCES class_subjects (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_enrollments_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE CASCADE;

ALTER TABLE grades ADD CONSTRAINT fk_grades_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_grades_enrollment_id FOREIGN KEY (enrollment_id) REFERENCES enrollments (id) ON DELETE CASCADE;

ALTER TABLE attendance ADD CONSTRAINT fk_attendance_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_attendance_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_attendance_schedule_id FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE CASCADE;

ALTER TABLE notifications ADD CONSTRAINT fk_notifications_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE fee_types ADD CONSTRAINT fk_fee_types_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE;

ALTER TABLE student_fees ADD CONSTRAINT fk_student_fees_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_student_fees_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_student_fees_fee_type_id FOREIGN KEY (fee_type_id) REFERENCES fee_types (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_student_fees_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE CASCADE;

-- ======================================================
-- PERFORMANCE INDEXES
-- ======================================================
-- Indexes for common queries
-- Tenant Users indexes
CREATE INDEX idx_tenant_users_tenant_id ON tenant_users (tenant_id);

CREATE INDEX idx_tenant_users_user_id ON tenant_users (user_id);

CREATE INDEX idx_tenant_users_is_active ON tenant_users (is_active);

CREATE INDEX idx_tenant_users_tenant_user ON tenant_users (tenant_id, user_id);

CREATE INDEX idx_tenant_users_user_active ON tenant_users (user_id, is_active);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_is_active ON users (is_active);

CREATE INDEX idx_users_full_name ON users (full_name);

CREATE INDEX idx_departments_tenant_id ON departments (tenant_id);

CREATE INDEX idx_departments_head_teacher_id ON departments (head_teacher_id);

CREATE INDEX idx_teachers_tenant_user_id ON teachers (tenant_user_id);

CREATE INDEX idx_teachers_department_id ON teachers (department_id);

CREATE INDEX idx_teachers_employee_number ON teachers (employee_number);

CREATE INDEX idx_parents_tenant_id ON parents (tenant_id);

CREATE INDEX idx_parents_phone ON parents (phone);

CREATE INDEX idx_parents_email ON parents (email);

CREATE INDEX idx_academic_years_tenant_id ON academic_years (tenant_id);

CREATE INDEX idx_academic_years_is_active ON academic_years (is_active);

CREATE INDEX idx_academic_years_dates ON academic_years (start_date, end_date);

CREATE INDEX idx_classes_tenant_id ON classes (tenant_id);

CREATE INDEX idx_classes_homeroom_teacher_id ON classes (homeroom_teacher_id);

CREATE INDEX idx_classes_academic_year_id ON classes (academic_year_id);

CREATE INDEX idx_classes_grade_level ON classes (grade_level);

CREATE INDEX idx_students_tenant_user_id ON students (tenant_user_id);

CREATE INDEX idx_students_student_number ON students (student_number);

CREATE INDEX idx_students_class_id ON students (class_id);

CREATE INDEX idx_students_parent_id ON students (parent_id);

CREATE INDEX idx_students_admission_date ON students (admission_date);

CREATE INDEX idx_subjects_tenant_id ON subjects (tenant_id);

CREATE INDEX idx_subjects_code ON subjects (code);

CREATE INDEX idx_subjects_department_id ON subjects (department_id);

CREATE INDEX idx_class_subjects_tenant_id ON class_subjects (tenant_id);

CREATE INDEX idx_class_subjects_class_id ON class_subjects (class_id);

CREATE INDEX idx_class_subjects_subject_id ON class_subjects (subject_id);

CREATE INDEX idx_class_subjects_teacher_id ON class_subjects (teacher_id);

CREATE INDEX idx_schedules_tenant_id ON schedules (tenant_id);

CREATE INDEX idx_schedules_class_subject_id ON schedules (class_subject_id);

CREATE INDEX idx_schedules_day_time ON schedules (day_of_week, start_time, end_time);

CREATE INDEX idx_schedules_room ON schedules (room);

CREATE INDEX idx_enrollments_tenant_id ON enrollments (tenant_id);

CREATE INDEX idx_enrollments_student_id ON enrollments (student_id);

CREATE INDEX idx_enrollments_class_subject_id ON enrollments (class_subject_id);

CREATE INDEX idx_enrollments_academic_year_id ON enrollments (academic_year_id);

CREATE INDEX idx_grades_tenant_id ON grades (tenant_id);

CREATE INDEX idx_grades_enrollment_id ON grades (enrollment_id);

CREATE INDEX idx_grades_grade_type ON grades (grade_type);

CREATE INDEX idx_grades_score ON grades (score);

CREATE INDEX idx_attendance_tenant_id ON attendance (tenant_id);

CREATE INDEX idx_attendance_student_id ON attendance (student_id);

CREATE INDEX idx_attendance_schedule_id ON attendance (schedule_id);

CREATE INDEX idx_attendance_date ON attendance (attendance_date);

CREATE INDEX idx_attendance_status ON attendance (status);

CREATE INDEX idx_attendance_student_date ON attendance (student_id, attendance_date);

CREATE INDEX idx_notifications_tenant_id ON notifications (tenant_id);

CREATE INDEX idx_notifications_user_id ON notifications (user_id);

CREATE INDEX idx_notifications_is_read ON notifications (is_read);

CREATE INDEX idx_notifications_user_unread ON notifications (user_id, is_read)
WHERE
  is_read = FALSE;

CREATE INDEX idx_fee_types_tenant_id ON fee_types (tenant_id);

CREATE INDEX idx_fee_types_name ON fee_types (name);

CREATE INDEX idx_fee_types_is_active ON fee_types (is_active);

CREATE INDEX idx_fee_types_is_mandatory ON fee_types (is_mandatory);

CREATE INDEX idx_student_fees_tenant_id ON student_fees (tenant_id);

CREATE INDEX idx_student_fees_student_id ON student_fees (student_id);

CREATE INDEX idx_student_fees_fee_type_id ON student_fees (fee_type_id);

CREATE INDEX idx_student_fees_academic_year_id ON student_fees (academic_year_id);

CREATE INDEX idx_student_fees_status ON student_fees (status);

CREATE INDEX idx_student_fees_due_date ON student_fees (due_date);

CREATE INDEX idx_student_fees_payment_date ON student_fees (payment_date);

CREATE INDEX idx_student_fees_student_status ON student_fees (student_id, status);

CREATE INDEX idx_student_fees_unpaid_overdue ON student_fees (due_date, status)
WHERE
  status IN ('unpaid', 'overdue');

-- Composite indexes for common query patterns
CREATE INDEX idx_attendance_student_date_status ON attendance (student_id, attendance_date, status);

CREATE INDEX idx_grades_enrollment_type ON grades (enrollment_id, grade_type);

CREATE INDEX idx_schedules_day_time_room ON schedules (day_of_week, start_time, room);

-- Additional performance indexes for common business queries

-- Subscription and billing performance indexes
CREATE INDEX idx_tenants_subscription_status ON tenants (subscription_status);
CREATE INDEX idx_tenants_plan_status ON tenants (plan_id, subscription_status);
CREATE INDEX idx_subscriptions_status_dates ON subscriptions (status, end_date);
CREATE INDEX idx_subscriptions_tenant_active ON subscriptions (tenant_id, status) WHERE status = 'active';
CREATE INDEX idx_invoices_status_due ON invoices (status, due_date);
CREATE INDEX idx_invoices_tenant_unpaid ON invoices (tenant_id, status) WHERE status IN ('unpaid', 'overdue');

-- User management and authentication indexes
CREATE INDEX idx_users_email_username ON users (email, username);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX idx_teachers_tenant_user_active ON teachers (tenant_user_id);
CREATE INDEX idx_students_tenant_user_active ON students (tenant_user_id);

-- Academic performance indexes
CREATE INDEX idx_classes_tenant_year ON classes (tenant_id, academic_year_id);
CREATE INDEX idx_class_subjects_class_teacher ON class_subjects (class_id, teacher_id);
CREATE INDEX idx_enrollments_student_year ON enrollments (student_id, academic_year_id);
CREATE INDEX idx_grades_student_type ON grades (enrollment_id, grade_type, score);
CREATE INDEX idx_attendance_class_date ON attendance (schedule_id, attendance_date);
CREATE INDEX idx_attendance_tenant_date ON attendance (tenant_id, attendance_date);

-- Fee management performance indexes
CREATE INDEX idx_fee_types_tenant_mandatory ON fee_types (tenant_id, is_mandatory, is_active);
CREATE INDEX idx_student_fees_overdue_report ON student_fees (tenant_id, status, due_date) WHERE status IN ('unpaid', 'overdue');
CREATE INDEX idx_student_fees_payment_history ON student_fees (student_id, payment_date) WHERE payment_date IS NOT NULL;
CREATE INDEX idx_student_fees_year_status ON student_fees (academic_year_id, status);

-- Department and organizational indexes
CREATE INDEX idx_departments_tenant_head ON departments (tenant_id, head_teacher_id);
CREATE INDEX idx_subjects_dept_active ON subjects (department_id, tenant_id);
CREATE INDEX idx_teachers_dept_active ON teachers (department_id, tenant_id);

-- Scheduling and timetable indexes
CREATE INDEX idx_schedules_room_day ON schedules (room, day_of_week);
CREATE INDEX idx_schedules_teacher_day ON schedules (class_subject_id, day_of_week);
CREATE INDEX idx_class_subjects_teacher_tenant ON class_subjects (teacher_id, tenant_id);

-- Communication and notification indexes
CREATE INDEX idx_notifications_tenant_unread ON notifications (tenant_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_user_recent ON notifications (user_id, created_at);

-- Parent and family relationship indexes
CREATE INDEX idx_parents_tenant_contact ON parents (tenant_id, phone);
CREATE INDEX idx_students_parent_tenant ON students (parent_id, tenant_id);

-- Feature flags and tenant capabilities
CREATE INDEX idx_tenant_features_enabled ON tenant_features (tenant_id, enabled) WHERE enabled = TRUE;
CREATE INDEX idx_feature_flags_code ON feature_flags (code);

-- Time-based performance indexes for reporting
CREATE INDEX idx_academic_years_tenant_active ON academic_years (tenant_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_students_admission_year ON students (tenant_id, admission_date);
CREATE INDEX idx_teachers_hire_date ON teachers (tenant_id, hire_date);

-- Full-text search preparation indexes (for future implementation)
CREATE INDEX idx_users_full_name_gin ON users USING gin(to_tsvector('english', full_name));
CREATE INDEX idx_subjects_name_gin ON subjects USING gin(to_tsvector('english', name));
CREATE INDEX idx_departments_name_gin ON departments USING gin(to_tsvector('english', name));

-- ======================================================
-- SAMPLE DATA
-- ======================================================
INSERT INTO
  roles (name, description)
VALUES
  ('Developer', 'Application Developer'),
  ('Admin', 'System Administrator'),
  ('Teacher', 'Teaching staff'),
  ('Student', 'Enrolled student'),
  ('Parent', 'Parent or guardian'),
  ('Staff', 'School staff');

-- =========================================
-- ENABLE ROW LEVEL SECURITY (RLS)
-- =========================================

-- Enable RLS per tenant on all tenant-based tables
ALTER TABLE tenant_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE teachers ENABLE ROW LEVEL SECURITY;
ALTER TABLE students ENABLE ROW LEVEL SECURITY;
ALTER TABLE subjects ENABLE ROW LEVEL SECURITY;
ALTER TABLE classes ENABLE ROW LEVEL SECURITY;
ALTER TABLE class_subjects ENABLE ROW LEVEL SECURITY;
ALTER TABLE academic_years ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance ENABLE ROW LEVEL SECURITY;
ALTER TABLE grades ENABLE ROW LEVEL SECURITY;
ALTER TABLE enrollments ENABLE ROW LEVEL SECURITY;
ALTER TABLE schedules ENABLE ROW LEVEL SECURITY;
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
ALTER TABLE parents ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_fees ENABLE ROW LEVEL SECURITY;
ALTER TABLE fee_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_features ENABLE ROW LEVEL SECURITY;

-- =========================================
-- POLICIES
-- =========================================
-- Assume we set the tenant context using: SET app.current_tenant = '<tenant_uuid>';

CREATE FUNCTION current_tenant_id() RETURNS UUID AS $$
BEGIN
    RETURN current_setting('app.current_tenant', true)::UUID;
END;
$$ LANGUAGE plpgsql STABLE;

-- Apply RLS Policy to restrict tenant access
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT tablename FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename IN ('tenant_users','teachers','students','subjects','classes','class_subjects','academic_years','attendance', 'grades','enrollments','schedules','departments','parents','notifications','student_fees','fee_types','tenant_features')
    LOOP
        EXECUTE format(
            'CREATE POLICY tenant_isolation ON %I
             USING (tenant_id = current_tenant_id())',
            tbl
        );
    END LOOP;
END;
$$;

-- =========================================
-- DEFAULT SETTINGS
-- =========================================
-- Example: before query, application sets tenant context:
-- SELECT set_config('app.current_tenant', '<tenant_uuid>', false);
