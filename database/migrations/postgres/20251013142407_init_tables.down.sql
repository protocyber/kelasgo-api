-- ======================================================
-- SCHOOL MANAGEMENT SYSTEM DATABASE SCHEMA - DOWN MIGRATION
-- PostgreSQL version
-- Author: Fitrah
-- ======================================================
-- This file undoes all changes made in the corresponding up migration

-- ======================================================
-- DROP RLS POLICIES
-- ======================================================
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT tablename FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename IN ('tenant_users','teachers','students','subjects','classes','class_subjects','academic_years','attendance', 'grades','enrollments','schedules','departments','parents','notifications','student_fees','fee_types','tenant_features')
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation ON %I', tbl);
    END LOOP;
END;
$$;

-- Drop RLS function
DROP FUNCTION IF EXISTS current_tenant_id();

-- ======================================================
-- DISABLE ROW LEVEL SECURITY
-- ======================================================
ALTER TABLE IF EXISTS tenant_features DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS fee_types DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS student_fees DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS notifications DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS parents DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS departments DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS schedules DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS enrollments DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS grades DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS attendance DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS academic_years DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS class_subjects DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS classes DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS subjects DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS students DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS teachers DISABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS tenant_users DISABLE ROW LEVEL SECURITY;

-- ======================================================
-- DROP INDEXES
-- ======================================================
-- Drop full-text search indexes
DROP INDEX IF EXISTS idx_departments_name_gin;
DROP INDEX IF EXISTS idx_subjects_name_gin;
DROP INDEX IF EXISTS idx_users_full_name_gin;

-- Drop time-based performance indexes for reporting
DROP INDEX IF EXISTS idx_teachers_hire_date;
DROP INDEX IF EXISTS idx_students_admission_year;
DROP INDEX IF EXISTS idx_academic_years_tenant_active;

-- Drop feature flags and tenant capabilities indexes
DROP INDEX IF EXISTS idx_feature_flags_code;
DROP INDEX IF EXISTS idx_tenant_features_enabled;

-- Drop parent and family relationship indexes
DROP INDEX IF EXISTS idx_students_parent_tenant;
DROP INDEX IF EXISTS idx_parents_tenant_contact;

-- Drop communication and notification indexes
DROP INDEX IF EXISTS idx_notifications_user_recent;
DROP INDEX IF EXISTS idx_notifications_tenant_unread;

-- Drop scheduling and timetable indexes
DROP INDEX IF EXISTS idx_class_subjects_teacher_tenant;
DROP INDEX IF EXISTS idx_schedules_teacher_day;
DROP INDEX IF EXISTS idx_schedules_room_day;

-- Drop department and organizational indexes
DROP INDEX IF EXISTS idx_teachers_dept_active;
DROP INDEX IF EXISTS idx_subjects_dept_active;
DROP INDEX IF EXISTS idx_departments_tenant_head;

-- Drop fee management performance indexes
DROP INDEX IF EXISTS idx_student_fees_year_status;
DROP INDEX IF EXISTS idx_student_fees_payment_history;
DROP INDEX IF EXISTS idx_student_fees_overdue_report;
DROP INDEX IF EXISTS idx_fee_types_tenant_mandatory;

-- Drop academic performance indexes
DROP INDEX IF EXISTS idx_attendance_tenant_date;
DROP INDEX IF EXISTS idx_attendance_class_date;
DROP INDEX IF EXISTS idx_grades_student_type;
DROP INDEX IF EXISTS idx_enrollments_student_year;
DROP INDEX IF EXISTS idx_class_subjects_class_teacher;
DROP INDEX IF EXISTS idx_classes_tenant_year;

-- Drop user management and authentication indexes
DROP INDEX IF EXISTS idx_students_tenant_user_active;
DROP INDEX IF EXISTS idx_teachers_tenant_user_active;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_users_email_username;

-- Drop subscription and billing performance indexes
DROP INDEX IF EXISTS idx_invoices_tenant_unpaid;
DROP INDEX IF EXISTS idx_invoices_status_due;
DROP INDEX IF EXISTS idx_subscriptions_tenant_active;
DROP INDEX IF EXISTS idx_subscriptions_status_dates;
DROP INDEX IF EXISTS idx_tenants_plan_status;
DROP INDEX IF EXISTS idx_tenants_subscription_status;

-- Drop composite indexes for common query patterns
DROP INDEX IF EXISTS idx_schedules_day_time_room;

DROP INDEX IF EXISTS idx_grades_enrollment_type;
DROP INDEX IF EXISTS idx_attendance_student_date_status;
DROP INDEX IF EXISTS idx_students_class_enrollment;


-- Drop partial indexes
DROP INDEX IF EXISTS idx_student_fees_unpaid_overdue;
DROP INDEX IF EXISTS idx_notifications_user_unread;

-- Drop all tenant_id indexes
DROP INDEX IF EXISTS idx_student_fees_tenant_id;
DROP INDEX IF EXISTS idx_fee_types_tenant_id;
DROP INDEX IF EXISTS idx_notifications_tenant_id;
DROP INDEX IF EXISTS idx_attendance_tenant_id;
DROP INDEX IF EXISTS idx_grades_tenant_id;
DROP INDEX IF EXISTS idx_enrollments_tenant_id;
DROP INDEX IF EXISTS idx_schedules_tenant_id;
DROP INDEX IF EXISTS idx_class_subjects_tenant_id;
DROP INDEX IF EXISTS idx_subjects_tenant_id;
DROP INDEX IF EXISTS idx_students_tenant_id;
DROP INDEX IF EXISTS idx_classes_tenant_id;
DROP INDEX IF EXISTS idx_academic_years_tenant_id;
DROP INDEX IF EXISTS idx_parents_tenant_id;
DROP INDEX IF EXISTS idx_teachers_tenant_id;
DROP INDEX IF EXISTS idx_departments_tenant_id;
-- Drop tenant_users indexes
DROP INDEX IF EXISTS idx_tenant_users_user_active;
DROP INDEX IF EXISTS idx_tenant_users_tenant_user;
DROP INDEX IF EXISTS idx_tenant_users_is_active;
DROP INDEX IF EXISTS idx_tenant_users_user_id;
DROP INDEX IF EXISTS idx_tenant_users_tenant_id;

-- Drop business logic indexes
DROP INDEX IF EXISTS idx_student_fees_student_status;
DROP INDEX IF EXISTS idx_student_fees_payment_date;
DROP INDEX IF EXISTS idx_student_fees_due_date;
DROP INDEX IF EXISTS idx_student_fees_status;
DROP INDEX IF EXISTS idx_student_fees_academic_year_id;
DROP INDEX IF EXISTS idx_student_fees_fee_type_id;
DROP INDEX IF EXISTS idx_student_fees_student_id;
DROP INDEX IF EXISTS idx_fee_types_is_mandatory;
DROP INDEX IF EXISTS idx_fee_types_is_active;
DROP INDEX IF EXISTS idx_fee_types_name;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_attendance_student_date;
DROP INDEX IF EXISTS idx_attendance_status;
DROP INDEX IF EXISTS idx_attendance_date;
DROP INDEX IF EXISTS idx_attendance_schedule_id;
DROP INDEX IF EXISTS idx_attendance_student_id;
DROP INDEX IF EXISTS idx_grades_score;
DROP INDEX IF EXISTS idx_grades_grade_type;
DROP INDEX IF EXISTS idx_grades_enrollment_id;
DROP INDEX IF EXISTS idx_enrollments_academic_year_id;
DROP INDEX IF EXISTS idx_enrollments_class_subject_id;
DROP INDEX IF EXISTS idx_enrollments_student_id;
DROP INDEX IF EXISTS idx_schedules_room;
DROP INDEX IF EXISTS idx_schedules_day_time;
DROP INDEX IF EXISTS idx_schedules_class_subject_id;
DROP INDEX IF EXISTS idx_class_subjects_teacher_id;
DROP INDEX IF EXISTS idx_class_subjects_subject_id;
DROP INDEX IF EXISTS idx_class_subjects_class_id;
DROP INDEX IF EXISTS idx_subjects_department_id;
DROP INDEX IF EXISTS idx_subjects_code;
DROP INDEX IF EXISTS idx_students_admission_date;
DROP INDEX IF EXISTS idx_students_parent_id;
DROP INDEX IF EXISTS idx_students_class_id;
DROP INDEX IF EXISTS idx_students_student_number;
DROP INDEX IF EXISTS idx_students_tenant_user_id;
DROP INDEX IF EXISTS idx_classes_grade_level;
DROP INDEX IF EXISTS idx_classes_academic_year_id;
DROP INDEX IF EXISTS idx_classes_homeroom_teacher_id;
DROP INDEX IF EXISTS idx_academic_years_dates;
DROP INDEX IF EXISTS idx_academic_years_is_active;
DROP INDEX IF EXISTS idx_parents_email;
DROP INDEX IF EXISTS idx_parents_phone;
DROP INDEX IF EXISTS idx_teachers_employee_number;
DROP INDEX IF EXISTS idx_teachers_department_id;
DROP INDEX IF EXISTS idx_teachers_tenant_user_id;
DROP INDEX IF EXISTS idx_departments_head_teacher_id;
DROP INDEX IF EXISTS idx_users_full_name;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_fee_types_updated_at;

DROP INDEX IF EXISTS idx_fee_types_created_at;

-- Drop notifications indexes
DROP INDEX IF EXISTS idx_notifications_is_read;

DROP INDEX IF EXISTS idx_notifications_user_id;

DROP INDEX IF EXISTS idx_notifications_updated_at;

DROP INDEX IF EXISTS idx_notifications_created_at;

-- Drop attendance indexes
DROP INDEX IF EXISTS idx_attendance_student_date;

DROP INDEX IF EXISTS idx_attendance_status;

DROP INDEX IF EXISTS idx_attendance_date;

DROP INDEX IF EXISTS idx_attendance_schedule_id;

DROP INDEX IF EXISTS idx_attendance_student_id;

DROP INDEX IF EXISTS idx_attendance_updated_at;

DROP INDEX IF EXISTS idx_attendance_created_at;

-- Drop grades indexes
DROP INDEX IF EXISTS idx_grades_score;

DROP INDEX IF EXISTS idx_grades_grade_type;

DROP INDEX IF EXISTS idx_grades_enrollment_id;

DROP INDEX IF EXISTS idx_grades_updated_at;

DROP INDEX IF EXISTS idx_grades_created_at;

-- Drop enrollments indexes
DROP INDEX IF EXISTS idx_enrollments_academic_year_id;

DROP INDEX IF EXISTS idx_enrollments_class_subject_id;

DROP INDEX IF EXISTS idx_enrollments_student_id;

DROP INDEX IF EXISTS idx_enrollments_updated_at;

DROP INDEX IF EXISTS idx_enrollments_created_at;

-- Drop schedules indexes
DROP INDEX IF EXISTS idx_schedules_room;

DROP INDEX IF EXISTS idx_schedules_day_time;

DROP INDEX IF EXISTS idx_schedules_class_subject_id;

DROP INDEX IF EXISTS idx_schedules_updated_at;

DROP INDEX IF EXISTS idx_schedules_created_at;

-- Drop class_subjects indexes
DROP INDEX IF EXISTS idx_class_subjects_teacher_id;

DROP INDEX IF EXISTS idx_class_subjects_subject_id;

DROP INDEX IF EXISTS idx_class_subjects_class_id;

DROP INDEX IF EXISTS idx_class_subjects_updated_at;

DROP INDEX IF EXISTS idx_class_subjects_created_at;

-- Drop subjects indexes
DROP INDEX IF EXISTS idx_subjects_department_id;

DROP INDEX IF EXISTS idx_subjects_code;

DROP INDEX IF EXISTS idx_subjects_updated_at;

DROP INDEX IF EXISTS idx_subjects_created_at;

-- Drop students indexes
DROP INDEX IF EXISTS idx_students_enrollment_date;

DROP INDEX IF EXISTS idx_students_parent_id;

DROP INDEX IF EXISTS idx_students_class_id;

DROP INDEX IF EXISTS idx_students_student_number;

DROP INDEX IF EXISTS idx_students_updated_at;

DROP INDEX IF EXISTS idx_students_created_at;

-- Drop classes indexes
DROP INDEX IF EXISTS idx_classes_grade_level;

DROP INDEX IF EXISTS idx_classes_academic_year_id;

DROP INDEX IF EXISTS idx_classes_homeroom_teacher_id;

DROP INDEX IF EXISTS idx_classes_updated_at;

DROP INDEX IF EXISTS idx_classes_created_at;

-- Drop academic_years indexes
DROP INDEX IF EXISTS idx_academic_years_dates;

DROP INDEX IF EXISTS idx_academic_years_is_active;

DROP INDEX IF EXISTS idx_academic_years_updated_at;

DROP INDEX IF EXISTS idx_academic_years_created_at;

-- Drop parents indexes
DROP INDEX IF EXISTS idx_parents_email;

DROP INDEX IF EXISTS idx_parents_phone;

DROP INDEX IF EXISTS idx_parents_updated_at;

DROP INDEX IF EXISTS idx_parents_created_at;

-- Drop teachers indexes
DROP INDEX IF EXISTS idx_teachers_employee_number;

DROP INDEX IF EXISTS idx_teachers_department_id;

DROP INDEX IF EXISTS idx_teachers_updated_at;

DROP INDEX IF EXISTS idx_teachers_created_at;

-- Drop departments indexes
DROP INDEX IF EXISTS idx_departments_head_teacher_id;

DROP INDEX IF EXISTS idx_departments_updated_at;

DROP INDEX IF EXISTS idx_departments_created_at;

-- Drop users indexes
DROP INDEX IF EXISTS idx_users_full_name;

DROP INDEX IF EXISTS idx_users_is_active;

DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_role_id;

DROP INDEX IF EXISTS idx_users_updated_at;

DROP INDEX IF EXISTS idx_users_created_at;

-- Drop roles indexes
DROP INDEX IF EXISTS idx_roles_updated_at;

DROP INDEX IF EXISTS idx_roles_created_at;

-- ======================================================
-- DROP FOREIGN KEY CONSTRAINTS FOR BUSINESS LOGIC
-- ======================================================
-- Drop business relationship constraints
ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_academic_year_id;

ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_fee_type_id;

ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_student_id;

ALTER TABLE notifications
DROP CONSTRAINT IF EXISTS fk_notifications_user_id;

ALTER TABLE attendance
DROP CONSTRAINT IF EXISTS fk_attendance_schedule_id;

ALTER TABLE attendance
DROP CONSTRAINT IF EXISTS fk_attendance_student_id;

ALTER TABLE grades
DROP CONSTRAINT IF EXISTS fk_grades_enrollment_id;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_academic_year_id;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_class_subject_id;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_student_id;

ALTER TABLE schedules
DROP CONSTRAINT IF EXISTS fk_schedules_class_subject_id;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_teacher_id;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_subject_id;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_class_id;

ALTER TABLE subjects
DROP CONSTRAINT IF EXISTS fk_subjects_department_id;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_parent_id;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_class_id;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_tenant_user_id;

ALTER TABLE classes
DROP CONSTRAINT IF EXISTS fk_classes_academic_year_id;

ALTER TABLE classes
DROP CONSTRAINT IF EXISTS fk_classes_homeroom_teacher_id;

ALTER TABLE departments
DROP CONSTRAINT IF EXISTS fk_departments_head_teacher_id;

ALTER TABLE teachers
DROP CONSTRAINT IF EXISTS fk_teachers_department_id;

ALTER TABLE teachers
DROP CONSTRAINT IF EXISTS fk_teachers_tenant_user_id;

-- Drop all tenant-related foreign key constraints
ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_tenant_id;

ALTER TABLE fee_types
DROP CONSTRAINT IF EXISTS fk_fee_types_tenant_id;

ALTER TABLE notifications
DROP CONSTRAINT IF EXISTS fk_notifications_tenant_id;

ALTER TABLE attendance
DROP CONSTRAINT IF EXISTS fk_attendance_tenant_id;

ALTER TABLE grades
DROP CONSTRAINT IF EXISTS fk_grades_tenant_id;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_tenant_id;

ALTER TABLE schedules
DROP CONSTRAINT IF EXISTS fk_schedules_tenant_id;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_tenant_id;

ALTER TABLE subjects
DROP CONSTRAINT IF EXISTS fk_subjects_tenant_id;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_tenant_id;

ALTER TABLE classes
DROP CONSTRAINT IF EXISTS fk_classes_tenant_id;

ALTER TABLE academic_years
DROP CONSTRAINT IF EXISTS fk_academic_years_tenant_id;

ALTER TABLE parents
DROP CONSTRAINT IF EXISTS fk_parents_tenant_id;

ALTER TABLE teachers
DROP CONSTRAINT IF EXISTS fk_teachers_tenant_id;

ALTER TABLE departments
DROP CONSTRAINT IF EXISTS fk_departments_tenant_id;

ALTER TABLE tenant_users
DROP CONSTRAINT IF EXISTS fk_tenant_users_tenant_id;

ALTER TABLE tenant_users
DROP CONSTRAINT IF EXISTS fk_tenant_users_user_id;

ALTER TABLE tenant_features
DROP CONSTRAINT IF EXISTS fk_tenant_features_tenant_id;

ALTER TABLE tenant_features
DROP CONSTRAINT IF EXISTS fk_tenant_features_feature_id;

-- Drop subscription and billing foreign key constraints
ALTER TABLE invoices
DROP CONSTRAINT IF EXISTS fk_invoices_tenant_id;

ALTER TABLE invoices
DROP CONSTRAINT IF EXISTS fk_invoices_subscription_id;

ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS fk_subscriptions_tenant_id;

ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS fk_subscriptions_plan_id;

ALTER TABLE tenants
DROP CONSTRAINT IF EXISTS fk_tenants_plan_id;

-- Drop user roles constraints
ALTER TABLE user_roles
DROP CONSTRAINT IF EXISTS fk_user_roles_user_id;

ALTER TABLE user_roles
DROP CONSTRAINT IF EXISTS fk_user_roles_role_id;

-- ======================================================
-- DROP FOREIGN KEY CONSTRAINTS FOR AUDIT COLUMNS
-- ======================================================
-- Drop audit constraints for all tables
ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_updated_by;

ALTER TABLE student_fees
DROP CONSTRAINT IF EXISTS fk_student_fees_created_by;

ALTER TABLE fee_types
DROP CONSTRAINT IF EXISTS fk_fee_types_updated_by;

ALTER TABLE fee_types
DROP CONSTRAINT IF EXISTS fk_fee_types_created_by;

ALTER TABLE notifications
DROP CONSTRAINT IF EXISTS fk_notifications_updated_by;

ALTER TABLE notifications
DROP CONSTRAINT IF EXISTS fk_notifications_created_by;

ALTER TABLE attendance
DROP CONSTRAINT IF EXISTS fk_attendance_updated_by;

ALTER TABLE attendance
DROP CONSTRAINT IF EXISTS fk_attendance_created_by;

ALTER TABLE grades
DROP CONSTRAINT IF EXISTS fk_grades_updated_by;

ALTER TABLE grades
DROP CONSTRAINT IF EXISTS fk_grades_created_by;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_updated_by;

ALTER TABLE enrollments
DROP CONSTRAINT IF EXISTS fk_enrollments_created_by;

ALTER TABLE schedules
DROP CONSTRAINT IF EXISTS fk_schedules_updated_by;

ALTER TABLE schedules
DROP CONSTRAINT IF EXISTS fk_schedules_created_by;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_updated_by;

ALTER TABLE class_subjects
DROP CONSTRAINT IF EXISTS fk_class_subjects_created_by;

ALTER TABLE subjects
DROP CONSTRAINT IF EXISTS fk_subjects_updated_by;

ALTER TABLE subjects
DROP CONSTRAINT IF EXISTS fk_subjects_created_by;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_updated_by;

ALTER TABLE students
DROP CONSTRAINT IF EXISTS fk_students_created_by;

ALTER TABLE classes
DROP CONSTRAINT IF EXISTS fk_classes_updated_by;

ALTER TABLE classes
DROP CONSTRAINT IF EXISTS fk_classes_created_by;

ALTER TABLE academic_years
DROP CONSTRAINT IF EXISTS fk_academic_years_updated_by;

ALTER TABLE academic_years
DROP CONSTRAINT IF EXISTS fk_academic_years_created_by;

ALTER TABLE parents
DROP CONSTRAINT IF EXISTS fk_parents_updated_by;

ALTER TABLE parents
DROP CONSTRAINT IF EXISTS fk_parents_created_by;

ALTER TABLE teachers
DROP CONSTRAINT IF EXISTS fk_teachers_updated_by;

ALTER TABLE teachers
DROP CONSTRAINT IF EXISTS fk_teachers_created_by;

ALTER TABLE departments
DROP CONSTRAINT IF EXISTS fk_departments_updated_by;

ALTER TABLE departments
DROP CONSTRAINT IF EXISTS fk_departments_created_by;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_updated_by;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_created_by;

ALTER TABLE roles
DROP CONSTRAINT IF EXISTS fk_roles_updated_by;

ALTER TABLE roles
DROP CONSTRAINT IF EXISTS fk_roles_created_by;

-- ======================================================
-- DROP TABLES
-- ======================================================
-- Drop tables in reverse dependency order (child tables first, then parents)

-- Drop billing and subscription tables
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS subscription_plans;

-- Drop school management tables  
DROP TABLE IF EXISTS student_fees;
DROP TABLE IF EXISTS fee_types;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS class_subjects;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS academic_years;
DROP TABLE IF EXISTS parents;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS tenant_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- Drop tenant and feature tables
DROP TABLE IF EXISTS tenant_features;
DROP TABLE IF EXISTS feature_flags;

-- ======================================================
-- DROP ENUMS
-- ======================================================
DROP TYPE IF EXISTS attendance_status_enum;
DROP TYPE IF EXISTS invoice_status_enum;
DROP TYPE IF EXISTS subscription_plan_status_enum;
DROP TYPE IF EXISTS subscription_status_enum;
DROP TYPE IF EXISTS fee_status_enum;

-- ======================================================
-- DROP EXTENSIONS
-- ======================================================
DROP EXTENSION IF EXISTS "uuid-ossp";

-- ======================================================
-- DROP SCHEMA (if desired - commented out for safety)
-- ======================================================
-- WARNING: Uncommenting this will drop the entire public schema
-- DROP SCHEMA IF EXISTS public CASCADE;
-- CREATE SCHEMA public;
