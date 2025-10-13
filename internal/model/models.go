package model

import (
	"time"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy null.Int  `json:"created_by,omitempty"`
	UpdatedBy null.Int  `json:"updated_by,omitempty"`
}

// Role represents the roles table
type Role struct {
	BaseModel
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Relationships
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

// User represents the users table
type User struct {
	BaseModel
	RoleID       null.Int  `json:"role_id,omitempty"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Email        string    `gorm:"size:100;uniqueIndex" json:"email"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	Gender       string    `gorm:"size:10;check:gender IN ('Male', 'Female')" json:"gender"`
	DateOfBirth  null.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Address      string    `gorm:"type:text" json:"address"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`

	// Relationships
	Role          *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Teacher       *Teacher       `gorm:"foreignKey:UserID" json:"teacher,omitempty"`
	Student       *Student       `gorm:"foreignKey:UserID" json:"student,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// Department represents the departments table
type Department struct {
	BaseModel
	Name          string   `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description   string   `gorm:"type:text" json:"description"`
	HeadTeacherID null.Int `json:"head_teacher_id,omitempty"`

	// Relationships
	HeadTeacher *Teacher  `gorm:"foreignKey:HeadTeacherID" json:"head_teacher,omitempty"`
	Teachers    []Teacher `gorm:"foreignKey:DepartmentID" json:"teachers,omitempty"`
	Subjects    []Subject `gorm:"foreignKey:DepartmentID" json:"subjects,omitempty"`
}

// Teacher represents the teachers table
type Teacher struct {
	BaseModel
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	EmployeeNumber string    `gorm:"size:50;uniqueIndex" json:"employee_number"`
	HireDate       null.Time `gorm:"type:date" json:"hire_date,omitempty"`
	DepartmentID   null.Int  `json:"department_id,omitempty"`
	Qualification  string    `gorm:"size:100" json:"qualification"`
	Position       string    `gorm:"size:100" json:"position"`

	// Relationships
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Department      *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Classes         []Class        `gorm:"foreignKey:HomeroomTeacherID" json:"classes,omitempty"`
	ClassSubjects   []ClassSubject `gorm:"foreignKey:TeacherID" json:"class_subjects,omitempty"`
	HeadDepartments []Department   `gorm:"foreignKey:HeadTeacherID" json:"head_departments,omitempty"`
}

// Parent represents the parents table
type Parent struct {
	BaseModel
	FullName     string `gorm:"size:100;not null" json:"full_name"`
	Phone        string `gorm:"size:20" json:"phone"`
	Email        string `gorm:"size:100" json:"email"`
	Address      string `gorm:"type:text" json:"address"`
	Relationship string `gorm:"size:50" json:"relationship"`

	// Relationships
	Students []Student `gorm:"foreignKey:ParentID" json:"students,omitempty"`
}

// AcademicYear represents the academic_years table
type AcademicYear struct {
	BaseModel
	Name      string    `gorm:"size:50;not null" json:"name"`
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`

	// Relationships
	Classes     []Class      `gorm:"foreignKey:AcademicYearID" json:"classes,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:AcademicYearID" json:"enrollments,omitempty"`
	StudentFees []StudentFee `gorm:"foreignKey:AcademicYearID" json:"student_fees,omitempty"`
}

// Class represents the classes table
type Class struct {
	BaseModel
	Name              string   `gorm:"size:50;not null" json:"name"`
	GradeLevel        int      `json:"grade_level"`
	HomeroomTeacherID null.Int `json:"homeroom_teacher_id,omitempty"`
	AcademicYearID    null.Int `json:"academic_year_id,omitempty"`

	// Relationships
	HomeroomTeacher *Teacher       `gorm:"foreignKey:HomeroomTeacherID" json:"homeroom_teacher,omitempty"`
	AcademicYear    *AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	Students        []Student      `gorm:"foreignKey:ClassID" json:"students,omitempty"`
	ClassSubjects   []ClassSubject `gorm:"foreignKey:ClassID" json:"class_subjects,omitempty"`
}

// Student represents the students table
type Student struct {
	BaseModel
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	StudentNumber  string    `gorm:"size:50;uniqueIndex;not null" json:"student_number"`
	EnrollmentDate time.Time `gorm:"type:date;not null" json:"enrollment_date"`
	ClassID        null.Int  `json:"class_id,omitempty"`
	ParentID       null.Int  `json:"parent_id,omitempty"`

	// Relationships
	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Class       *Class       `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Parent      *Parent      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:StudentID" json:"enrollments,omitempty"`
	Attendance  []Attendance `gorm:"foreignKey:StudentID" json:"attendance,omitempty"`
	StudentFees []StudentFee `gorm:"foreignKey:StudentID" json:"student_fees,omitempty"`
}

// Subject represents the subjects table
type Subject struct {
	BaseModel
	Name         string   `gorm:"size:100;not null" json:"name"`
	Code         string   `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description  string   `gorm:"type:text" json:"description"`
	DepartmentID null.Int `json:"department_id,omitempty"`
	Credit       int      `gorm:"default:0" json:"credit"`

	// Relationships
	Department    *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	ClassSubjects []ClassSubject `gorm:"foreignKey:SubjectID" json:"class_subjects,omitempty"`
}

// ClassSubject represents the class_subjects table (linking class, subject, teacher)
type ClassSubject struct {
	BaseModel
	ClassID   null.Int `json:"class_id,omitempty"`
	SubjectID null.Int `json:"subject_id,omitempty"`
	TeacherID null.Int `json:"teacher_id,omitempty"`

	// Relationships
	Class       *Class       `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Subject     *Subject     `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher     *Teacher     `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Schedules   []Schedule   `gorm:"foreignKey:ClassSubjectID" json:"schedules,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:ClassSubjectID" json:"enrollments,omitempty"`
}

// Schedule represents the schedules table
type Schedule struct {
	BaseModel
	ClassSubjectID null.Int `json:"class_subject_id,omitempty"`
	DayOfWeek      string   `gorm:"size:15;check:day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday')" json:"day_of_week"`
	StartTime      string   `gorm:"type:time" json:"start_time"`
	EndTime        string   `gorm:"type:time" json:"end_time"`
	Room           string   `gorm:"size:50" json:"room"`

	// Relationships
	ClassSubject *ClassSubject `gorm:"foreignKey:ClassSubjectID" json:"class_subject,omitempty"`
	Attendance   []Attendance  `gorm:"foreignKey:ScheduleID" json:"attendance,omitempty"`
}

// Enrollment represents the enrollments table
type Enrollment struct {
	BaseModel
	StudentID      null.Int `json:"student_id,omitempty"`
	ClassSubjectID null.Int `json:"class_subject_id,omitempty"`
	AcademicYearID null.Int `json:"academic_year_id,omitempty"`

	// Relationships
	Student      *Student      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	ClassSubject *ClassSubject `gorm:"foreignKey:ClassSubjectID" json:"class_subject,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	Grades       []Grade       `gorm:"foreignKey:EnrollmentID" json:"grades,omitempty"`
}

// Grade represents the grades table
type Grade struct {
	BaseModel
	EnrollmentID null.Int   `json:"enrollment_id,omitempty"`
	GradeType    string     `gorm:"size:50;check:grade_type IN ('Assignment', 'Midterm', 'Final', 'Other')" json:"grade_type"`
	Score        null.Float `gorm:"type:decimal(5,2)" json:"score,omitempty"`
	Remarks      string     `gorm:"type:text" json:"remarks"`

	// Relationships
	Enrollment *Enrollment `gorm:"foreignKey:EnrollmentID" json:"enrollment,omitempty"`
}

// Attendance represents the attendance table
type Attendance struct {
	BaseModel
	StudentID      null.Int  `json:"student_id,omitempty"`
	ScheduleID     null.Int  `json:"schedule_id,omitempty"`
	Status         string    `gorm:"size:20;check:status IN ('Present', 'Absent', 'Late', 'Excused')" json:"status"`
	AttendanceDate time.Time `gorm:"type:date;default:CURRENT_DATE" json:"attendance_date"`
	Remarks        string    `gorm:"type:text" json:"remarks"`

	// Relationships
	Student  *Student  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Schedule *Schedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// Notification represents the notifications table
type Notification struct {
	BaseModel
	UserID  null.Int `json:"user_id,omitempty"`
	Title   string   `gorm:"size:100" json:"title"`
	Message string   `gorm:"type:text" json:"message"`
	IsRead  bool     `gorm:"default:false" json:"is_read"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// FeeType represents the fee_types table
type FeeType struct {
	BaseModel
	Name          string     `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	DefaultAmount null.Float `gorm:"type:decimal(10,2);default:0;check:default_amount >= 0" json:"default_amount,omitempty"`
	IsMandatory   bool       `gorm:"default:true" json:"is_mandatory"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`

	// Relationships
	StudentFees []StudentFee `gorm:"foreignKey:FeeTypeID" json:"student_fees,omitempty"`
}

// FeeStatus represents the fee status enum
type FeeStatus string

const (
	FeeStatusPaid    FeeStatus = "paid"
	FeeStatusUnpaid  FeeStatus = "unpaid"
	FeeStatusPartial FeeStatus = "partial"
	FeeStatusOverdue FeeStatus = "overdue"
)

// StudentFee represents the student_fees table
type StudentFee struct {
	BaseModel
	StudentID      null.Int  `json:"student_id,omitempty"`
	FeeTypeID      null.Int  `json:"fee_type_id,omitempty"`
	AcademicYearID null.Int  `json:"academic_year_id,omitempty"`
	Amount         float64   `gorm:"type:decimal(10,2);not null;check:amount >= 0" json:"amount"`
	DueDate        time.Time `gorm:"type:date;not null" json:"due_date"`
	Status         FeeStatus `gorm:"type:fee_status_enum;default:'unpaid'" json:"status"`
	PaymentDate    null.Time `gorm:"type:date" json:"payment_date,omitempty"`
	PaymentMethod  string    `gorm:"size:50" json:"payment_method"`
	Notes          string    `gorm:"type:text" json:"notes"`

	// Relationships
	Student      *Student      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	FeeType      *FeeType      `gorm:"foreignKey:FeeTypeID" json:"fee_type,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
}

// BeforeCreate hook to set audit fields
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to update audit fields
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
