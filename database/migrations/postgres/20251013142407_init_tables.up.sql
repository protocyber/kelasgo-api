-- ======================================================
-- SCHOOL MANAGEMENT SYSTEM DATABASE SCHEMA
-- PostgreSQL version
-- Author: Fitrah
-- ======================================================
-- ======================================================
-- ENUMS
-- ======================================================
CREATE TYPE fee_status_enum AS ENUM ('paid', 'unpaid', 'partial', 'overdue');

-- ======================================================
-- 1. ROLES
-- ======================================================
CREATE TABLE
  roles (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
  );

-- ======================================================
-- 2. USERS
-- ======================================================
CREATE TABLE
  users (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    role_id INT,
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
-- 3. DEPARTMENTS
-- ======================================================
CREATE TABLE
  departments (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    head_teacher_id INT
  );

-- ======================================================
-- 4. TEACHERS
-- ======================================================
CREATE TABLE
  teachers (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    employee_number VARCHAR(50) UNIQUE,
    hire_date DATE,
    department_id INT,
    qualification VARCHAR(100),
    position VARCHAR(100)
  );

-- ======================================================
-- 5. PARENTS
-- ======================================================
CREATE TABLE
  parents (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    relationship VARCHAR(50)
  );

-- ======================================================
-- 6. SCHOOL YEARS
-- ======================================================
CREATE TABLE
  academic_years (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE
  );

-- ======================================================
-- 7. CLASSES
-- ======================================================
CREATE TABLE
  classes (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    grade_level INT,
    homeroom_teacher_id INT,
    academic_year_id INT
  );

-- ======================================================
-- 8. STUDENTS
-- ======================================================
CREATE TABLE
  students (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    student_number VARCHAR(50) UNIQUE NOT NULL,
    enrollment_date DATE NOT NULL,
    class_id INT,
    parent_id INT
  );

-- ======================================================
-- 9. SUBJECTS
-- ======================================================
CREATE TABLE
  subjects (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    department_id INT,
    credit INT DEFAULT 0
  );

-- ======================================================
-- 10. CLASS SUBJECTS (link between class, subject, teacher)
-- ======================================================
CREATE TABLE
  class_subjects (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    class_id INT,
    subject_id INT,
    teacher_id INT,
    UNIQUE (class_id, subject_id, teacher_id)
  );

-- ======================================================
-- 11. SCHEDULES
-- ======================================================
CREATE TABLE
  schedules (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    class_subject_id INT,
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
-- 12. ENROLLMENTS
-- ======================================================
CREATE TABLE
  enrollments (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    student_id INT,
    class_subject_id INT,
    academic_year_id INT
  );

-- ======================================================
-- 13. GRADES
-- ======================================================
CREATE TABLE
  grades (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    enrollment_id INT,
    grade_type VARCHAR(50) CHECK (
      grade_type IN ('Assignment', 'Midterm', 'Final', 'Other')
    ),
    score DECIMAL(5, 2),
    remarks TEXT
  );

-- ======================================================
-- 14. ATTENDANCE
-- ======================================================
CREATE TABLE
  attendance (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    student_id INT,
    schedule_id INT,
    status VARCHAR(20) CHECK (
      status IN ('Present', 'Absent', 'Late', 'Excused')
    ),
    attendance_date DATE DEFAULT CURRENT_DATE,
    remarks TEXT
  );

-- ======================================================
-- 15. NOTIFICATIONS
-- ======================================================
CREATE TABLE
  notifications (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id SERIAL PRIMARY KEY,
    user_id INT,
    title VARCHAR(100),
    message TEXT,
    is_read BOOLEAN DEFAULT FALSE
  );

-- ========================================
-- TABLE: fee_types
-- ========================================
CREATE TABLE
  fee_types (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    default_amount DECIMAL(10, 2) DEFAULT 0,
    is_mandatory BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT chk_fee_types_amount CHECK (default_amount >= 0)
  );

-- ========================================
-- TABLE: student_fees
-- ========================================
CREATE TABLE
  student_fees (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    created_by INT,
    updated_by INT,
    id BIGSERIAL PRIMARY KEY,
    student_id INT,
    fee_type_id BIGINT,
    academic_year_id INT,
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
-- FOREIGN KEY CONSTRAINTS FOR AUDIT COLUMNS
-- ======================================================
-- Add foreign key constraints for created_by and updated_by columns
-- These reference the users table for audit tracking
ALTER TABLE roles ADD CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_roles_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE users ADD CONSTRAINT fk_users_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE departments ADD CONSTRAINT fk_departments_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_departments_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE teachers ADD CONSTRAINT fk_teachers_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_teachers_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE parents ADD CONSTRAINT fk_parents_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_parents_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE academic_years ADD CONSTRAINT fk_academic_years_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_academic_years_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE classes ADD CONSTRAINT fk_classes_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_classes_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE students ADD CONSTRAINT fk_students_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_students_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE subjects ADD CONSTRAINT fk_subjects_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_subjects_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE class_subjects ADD CONSTRAINT fk_class_subjects_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_class_subjects_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE schedules ADD CONSTRAINT fk_schedules_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_schedules_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE enrollments ADD CONSTRAINT fk_enrollments_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_enrollments_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE grades ADD CONSTRAINT fk_grades_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_grades_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE attendance ADD CONSTRAINT fk_attendance_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_attendance_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE notifications ADD CONSTRAINT fk_notifications_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_notifications_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE fee_types ADD CONSTRAINT fk_fee_types_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_fee_types_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

ALTER TABLE student_fees ADD CONSTRAINT fk_student_fees_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_student_fees_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;

-- ======================================================
-- FOREIGN KEY CONSTRAINTS FOR BUSINESS LOGIC
-- ======================================================
-- Add foreign key constraints for business relationships
-- Moved outside table definitions to avoid circular dependencies
ALTER TABLE users ADD CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE SET NULL;

ALTER TABLE teachers ADD CONSTRAINT fk_teachers_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_teachers_department_id FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE SET NULL;

ALTER TABLE departments ADD CONSTRAINT fk_departments_head_teacher_id FOREIGN KEY (head_teacher_id) REFERENCES teachers (id) ON DELETE SET NULL;

ALTER TABLE classes ADD CONSTRAINT fk_classes_homeroom_teacher_id FOREIGN KEY (homeroom_teacher_id) REFERENCES teachers (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_classes_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE SET NULL;

ALTER TABLE students ADD CONSTRAINT fk_students_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_students_class_id FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE SET NULL,
ADD CONSTRAINT fk_students_parent_id FOREIGN KEY (parent_id) REFERENCES parents (id) ON DELETE SET NULL;

ALTER TABLE subjects ADD CONSTRAINT fk_subjects_department_id FOREIGN KEY (department_id) REFERENCES departments (id) ON DELETE SET NULL;

ALTER TABLE class_subjects ADD CONSTRAINT fk_class_subjects_class_id FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_class_subjects_subject_id FOREIGN KEY (subject_id) REFERENCES subjects (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_class_subjects_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers (id) ON DELETE SET NULL;

ALTER TABLE schedules ADD CONSTRAINT fk_schedules_class_subject_id FOREIGN KEY (class_subject_id) REFERENCES class_subjects (id) ON DELETE CASCADE;

ALTER TABLE enrollments ADD CONSTRAINT fk_enrollments_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_enrollments_class_subject_id FOREIGN KEY (class_subject_id) REFERENCES class_subjects (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_enrollments_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE CASCADE;

ALTER TABLE grades ADD CONSTRAINT fk_grades_enrollment_id FOREIGN KEY (enrollment_id) REFERENCES enrollments (id) ON DELETE CASCADE;

ALTER TABLE attendance ADD CONSTRAINT fk_attendance_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_attendance_schedule_id FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE CASCADE;

ALTER TABLE notifications ADD CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE student_fees ADD CONSTRAINT fk_student_fees_student_id FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_student_fees_fee_type_id FOREIGN KEY (fee_type_id) REFERENCES fee_types (id) ON DELETE CASCADE,
ADD CONSTRAINT fk_student_fees_academic_year_id FOREIGN KEY (academic_year_id) REFERENCES academic_years (id) ON DELETE CASCADE;

-- ======================================================
-- PERFORMANCE INDEXES
-- ======================================================
-- Indexes for audit columns (common queries)
CREATE INDEX idx_roles_created_at ON roles (created_at);

CREATE INDEX idx_roles_updated_at ON roles (updated_at);

CREATE INDEX idx_users_created_at ON users (created_at);

CREATE INDEX idx_users_updated_at ON users (updated_at);

CREATE INDEX idx_users_role_id ON users (role_id);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_is_active ON users (is_active);

CREATE INDEX idx_users_full_name ON users (full_name);

CREATE INDEX idx_departments_created_at ON departments (created_at);

CREATE INDEX idx_departments_updated_at ON departments (updated_at);

CREATE INDEX idx_departments_head_teacher_id ON departments (head_teacher_id);

CREATE INDEX idx_teachers_created_at ON teachers (created_at);

CREATE INDEX idx_teachers_updated_at ON teachers (updated_at);

CREATE INDEX idx_teachers_user_id ON teachers (user_id);

CREATE INDEX idx_teachers_department_id ON teachers (department_id);

CREATE INDEX idx_teachers_employee_number ON teachers (employee_number);

CREATE INDEX idx_parents_created_at ON parents (created_at);

CREATE INDEX idx_parents_updated_at ON parents (updated_at);

CREATE INDEX idx_parents_phone ON parents (phone);

CREATE INDEX idx_parents_email ON parents (email);

CREATE INDEX idx_academic_years_created_at ON academic_years (created_at);

CREATE INDEX idx_academic_years_updated_at ON academic_years (updated_at);

CREATE INDEX idx_academic_years_is_active ON academic_years (is_active);

CREATE INDEX idx_academic_years_dates ON academic_years (start_date, end_date);

CREATE INDEX idx_classes_created_at ON classes (created_at);

CREATE INDEX idx_classes_updated_at ON classes (updated_at);

CREATE INDEX idx_classes_homeroom_teacher_id ON classes (homeroom_teacher_id);

CREATE INDEX idx_classes_academic_year_id ON classes (academic_year_id);

CREATE INDEX idx_classes_grade_level ON classes (grade_level);

CREATE INDEX idx_students_created_at ON students (created_at);

CREATE INDEX idx_students_updated_at ON students (updated_at);

CREATE INDEX idx_students_user_id ON students (user_id);

CREATE INDEX idx_students_student_number ON students (student_number);

CREATE INDEX idx_students_class_id ON students (class_id);

CREATE INDEX idx_students_parent_id ON students (parent_id);

CREATE INDEX idx_students_enrollment_date ON students (enrollment_date);

CREATE INDEX idx_subjects_created_at ON subjects (created_at);

CREATE INDEX idx_subjects_updated_at ON subjects (updated_at);

CREATE INDEX idx_subjects_code ON subjects (code);

CREATE INDEX idx_subjects_department_id ON subjects (department_id);

CREATE INDEX idx_class_subjects_created_at ON class_subjects (created_at);

CREATE INDEX idx_class_subjects_updated_at ON class_subjects (updated_at);

CREATE INDEX idx_class_subjects_class_id ON class_subjects (class_id);

CREATE INDEX idx_class_subjects_subject_id ON class_subjects (subject_id);

CREATE INDEX idx_class_subjects_teacher_id ON class_subjects (teacher_id);

CREATE INDEX idx_schedules_created_at ON schedules (created_at);

CREATE INDEX idx_schedules_updated_at ON schedules (updated_at);

CREATE INDEX idx_schedules_class_subject_id ON schedules (class_subject_id);

CREATE INDEX idx_schedules_day_time ON schedules (day_of_week, start_time, end_time);

CREATE INDEX idx_schedules_room ON schedules (room);

CREATE INDEX idx_enrollments_created_at ON enrollments (created_at);

CREATE INDEX idx_enrollments_updated_at ON enrollments (updated_at);

CREATE INDEX idx_enrollments_student_id ON enrollments (student_id);

CREATE INDEX idx_enrollments_class_subject_id ON enrollments (class_subject_id);

CREATE INDEX idx_enrollments_academic_year_id ON enrollments (academic_year_id);

CREATE INDEX idx_grades_created_at ON grades (created_at);

CREATE INDEX idx_grades_updated_at ON grades (updated_at);

CREATE INDEX idx_grades_enrollment_id ON grades (enrollment_id);

CREATE INDEX idx_grades_grade_type ON grades (grade_type);

CREATE INDEX idx_grades_score ON grades (score);

CREATE INDEX idx_attendance_created_at ON attendance (created_at);

CREATE INDEX idx_attendance_updated_at ON attendance (updated_at);

CREATE INDEX idx_attendance_student_id ON attendance (student_id);

CREATE INDEX idx_attendance_schedule_id ON attendance (schedule_id);

CREATE INDEX idx_attendance_date ON attendance (attendance_date);

CREATE INDEX idx_attendance_status ON attendance (status);

CREATE INDEX idx_attendance_student_date ON attendance (student_id, attendance_date);

CREATE INDEX idx_notifications_created_at ON notifications (created_at);

CREATE INDEX idx_notifications_updated_at ON notifications (updated_at);

CREATE INDEX idx_notifications_user_id ON notifications (user_id);

CREATE INDEX idx_notifications_is_read ON notifications (is_read);

CREATE INDEX idx_notifications_user_unread ON notifications (user_id, is_read)
WHERE
  is_read = FALSE;

CREATE INDEX idx_fee_types_created_at ON fee_types (created_at);

CREATE INDEX idx_fee_types_updated_at ON fee_types (updated_at);

CREATE INDEX idx_fee_types_name ON fee_types (name);

CREATE INDEX idx_fee_types_is_active ON fee_types (is_active);

CREATE INDEX idx_fee_types_is_mandatory ON fee_types (is_mandatory);

CREATE INDEX idx_student_fees_created_at ON student_fees (created_at);

CREATE INDEX idx_student_fees_updated_at ON student_fees (updated_at);

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
CREATE INDEX idx_users_role_active ON users (role_id, is_active);

CREATE INDEX idx_students_class_enrollment ON students (class_id, enrollment_date);

CREATE INDEX idx_attendance_student_date_status ON attendance (student_id, attendance_date, status);

CREATE INDEX idx_grades_enrollment_type ON grades (enrollment_id, grade_type);

CREATE INDEX idx_schedules_day_time_room ON schedules (day_of_week, start_time, room);

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
