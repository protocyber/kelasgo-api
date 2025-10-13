package dto

import (
	"time"

	"github.com/guregu/null/v5"
)

// Common response structures
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// Auth DTOs
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// User DTOs
type CreateUserRequest struct {
	RoleID      null.Int  `json:"role_id" validate:"omitempty,min=1"`
	Username    string    `json:"username" validate:"required,min=3,max=50"`
	Password    string    `json:"password" validate:"required,min=6"`
	Email       string    `json:"email" validate:"omitempty,email,max=100"`
	FullName    string    `json:"full_name" validate:"required,max=100"`
	Gender      string    `json:"gender" validate:"omitempty,oneof=Male Female"`
	DateOfBirth null.Time `json:"date_of_birth,omitempty"`
	Phone       string    `json:"phone" validate:"omitempty,max=20"`
	Address     string    `json:"address"`
	IsActive    *bool     `json:"is_active,omitempty"`
}

type UpdateUserRequest struct {
	RoleID      null.Int  `json:"role_id" validate:"omitempty,min=1"`
	Email       string    `json:"email" validate:"omitempty,email,max=100"`
	FullName    string    `json:"full_name" validate:"omitempty,max=100"`
	Gender      string    `json:"gender" validate:"omitempty,oneof=Male Female"`
	DateOfBirth null.Time `json:"date_of_birth,omitempty"`
	Phone       string    `json:"phone" validate:"omitempty,max=20"`
	Address     string    `json:"address"`
	IsActive    *bool     `json:"is_active,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// Role DTOs
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"omitempty,max=50"`
	Description string `json:"description"`
}

// Department DTOs
type CreateDepartmentRequest struct {
	Name          string   `json:"name" validate:"required,max=100"`
	Description   string   `json:"description"`
	HeadTeacherID null.Int `json:"head_teacher_id" validate:"omitempty,min=1"`
}

type UpdateDepartmentRequest struct {
	Name          string   `json:"name" validate:"omitempty,max=100"`
	Description   string   `json:"description"`
	HeadTeacherID null.Int `json:"head_teacher_id" validate:"omitempty,min=1"`
}

// Teacher DTOs
type CreateTeacherRequest struct {
	UserID         uint      `json:"user_id" validate:"required,min=1"`
	EmployeeNumber string    `json:"employee_number" validate:"omitempty,max=50"`
	HireDate       null.Time `json:"hire_date,omitempty"`
	DepartmentID   null.Int  `json:"department_id" validate:"omitempty,min=1"`
	Qualification  string    `json:"qualification" validate:"omitempty,max=100"`
	Position       string    `json:"position" validate:"omitempty,max=100"`
}

type UpdateTeacherRequest struct {
	EmployeeNumber string    `json:"employee_number" validate:"omitempty,max=50"`
	HireDate       null.Time `json:"hire_date,omitempty"`
	DepartmentID   null.Int  `json:"department_id" validate:"omitempty,min=1"`
	Qualification  string    `json:"qualification" validate:"omitempty,max=100"`
	Position       string    `json:"position" validate:"omitempty,max=100"`
}

// Student DTOs
type CreateStudentRequest struct {
	UserID         uint      `json:"user_id" validate:"required,min=1"`
	StudentNumber  string    `json:"student_number" validate:"required,max=50"`
	EnrollmentDate time.Time `json:"enrollment_date" validate:"required"`
	ClassID        null.Int  `json:"class_id" validate:"omitempty,min=1"`
	ParentID       null.Int  `json:"parent_id" validate:"omitempty,min=1"`
}

type UpdateStudentRequest struct {
	StudentNumber  string    `json:"student_number" validate:"omitempty,max=50"`
	EnrollmentDate time.Time `json:"enrollment_date,omitempty"`
	ClassID        null.Int  `json:"class_id" validate:"omitempty,min=1"`
	ParentID       null.Int  `json:"parent_id" validate:"omitempty,min=1"`
}

// Parent DTOs
type CreateParentRequest struct {
	FullName     string `json:"full_name" validate:"required,max=100"`
	Phone        string `json:"phone" validate:"omitempty,max=20"`
	Email        string `json:"email" validate:"omitempty,email,max=100"`
	Address      string `json:"address"`
	Relationship string `json:"relationship" validate:"omitempty,max=50"`
}

type UpdateParentRequest struct {
	FullName     string `json:"full_name" validate:"omitempty,max=100"`
	Phone        string `json:"phone" validate:"omitempty,max=20"`
	Email        string `json:"email" validate:"omitempty,email,max=100"`
	Address      string `json:"address"`
	Relationship string `json:"relationship" validate:"omitempty,max=50"`
}

// Academic Year DTOs
type CreateAcademicYearRequest struct {
	Name      string    `json:"name" validate:"required,max=50"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	IsActive  *bool     `json:"is_active,omitempty"`
}

type UpdateAcademicYearRequest struct {
	Name      string    `json:"name" validate:"omitempty,max=50"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
}

// Class DTOs
type CreateClassRequest struct {
	Name              string   `json:"name" validate:"required,max=50"`
	GradeLevel        int      `json:"grade_level" validate:"omitempty,min=1,max=12"`
	HomeroomTeacherID null.Int `json:"homeroom_teacher_id" validate:"omitempty,min=1"`
	AcademicYearID    null.Int `json:"academic_year_id" validate:"omitempty,min=1"`
}

type UpdateClassRequest struct {
	Name              string   `json:"name" validate:"omitempty,max=50"`
	GradeLevel        int      `json:"grade_level" validate:"omitempty,min=1,max=12"`
	HomeroomTeacherID null.Int `json:"homeroom_teacher_id" validate:"omitempty,min=1"`
	AcademicYearID    null.Int `json:"academic_year_id" validate:"omitempty,min=1"`
}

// Subject DTOs
type CreateSubjectRequest struct {
	Name         string   `json:"name" validate:"required,max=100"`
	Code         string   `json:"code" validate:"required,max=50"`
	Description  string   `json:"description"`
	DepartmentID null.Int `json:"department_id" validate:"omitempty,min=1"`
	Credit       int      `json:"credit" validate:"omitempty,min=0"`
}

type UpdateSubjectRequest struct {
	Name         string   `json:"name" validate:"omitempty,max=100"`
	Code         string   `json:"code" validate:"omitempty,max=50"`
	Description  string   `json:"description"`
	DepartmentID null.Int `json:"department_id" validate:"omitempty,min=1"`
	Credit       int      `json:"credit" validate:"omitempty,min=0"`
}

// Schedule DTOs
type CreateScheduleRequest struct {
	ClassSubjectID null.Int `json:"class_subject_id" validate:"omitempty,min=1"`
	DayOfWeek      string   `json:"day_of_week" validate:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday"`
	StartTime      string   `json:"start_time" validate:"required"`
	EndTime        string   `json:"end_time" validate:"required"`
	Room           string   `json:"room" validate:"omitempty,max=50"`
}

type UpdateScheduleRequest struct {
	ClassSubjectID null.Int `json:"class_subject_id" validate:"omitempty,min=1"`
	DayOfWeek      string   `json:"day_of_week" validate:"omitempty,oneof=Monday Tuesday Wednesday Thursday Friday Saturday"`
	StartTime      string   `json:"start_time,omitempty"`
	EndTime        string   `json:"end_time,omitempty"`
	Room           string   `json:"room" validate:"omitempty,max=50"`
}

// Grade DTOs
type CreateGradeRequest struct {
	EnrollmentID null.Int   `json:"enrollment_id" validate:"omitempty,min=1"`
	GradeType    string     `json:"grade_type" validate:"required,oneof=Assignment Midterm Final Other"`
	Score        null.Float `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks      string     `json:"remarks"`
}

type UpdateGradeRequest struct {
	GradeType string     `json:"grade_type" validate:"omitempty,oneof=Assignment Midterm Final Other"`
	Score     null.Float `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks   string     `json:"remarks"`
}

// Attendance DTOs
type CreateAttendanceRequest struct {
	StudentID      null.Int  `json:"student_id" validate:"omitempty,min=1"`
	ScheduleID     null.Int  `json:"schedule_id" validate:"omitempty,min=1"`
	Status         string    `json:"status" validate:"required,oneof=Present Absent Late Excused"`
	AttendanceDate time.Time `json:"attendance_date,omitempty"`
	Remarks        string    `json:"remarks"`
}

type UpdateAttendanceRequest struct {
	Status         string    `json:"status" validate:"omitempty,oneof=Present Absent Late Excused"`
	AttendanceDate time.Time `json:"attendance_date,omitempty"`
	Remarks        string    `json:"remarks"`
}

// Fee Type DTOs
type CreateFeeTypeRequest struct {
	Name          string     `json:"name" validate:"required,max=100"`
	Description   string     `json:"description"`
	DefaultAmount null.Float `json:"default_amount,omitempty" validate:"omitempty,min=0"`
	IsMandatory   *bool      `json:"is_mandatory,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type UpdateFeeTypeRequest struct {
	Name          string     `json:"name" validate:"omitempty,max=100"`
	Description   string     `json:"description"`
	DefaultAmount null.Float `json:"default_amount,omitempty" validate:"omitempty,min=0"`
	IsMandatory   *bool      `json:"is_mandatory,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

// Student Fee DTOs
type CreateStudentFeeRequest struct {
	StudentID      null.Int  `json:"student_id" validate:"omitempty,min=1"`
	FeeTypeID      null.Int  `json:"fee_type_id" validate:"omitempty,min=1"`
	AcademicYearID null.Int  `json:"academic_year_id" validate:"omitempty,min=1"`
	Amount         float64   `json:"amount" validate:"required,min=0"`
	DueDate        time.Time `json:"due_date" validate:"required"`
	Status         string    `json:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
	PaymentDate    null.Time `json:"payment_date,omitempty"`
	PaymentMethod  string    `json:"payment_method" validate:"omitempty,max=50"`
	Notes          string    `json:"notes"`
}

type UpdateStudentFeeRequest struct {
	Amount        float64   `json:"amount" validate:"omitempty,min=0"`
	DueDate       time.Time `json:"due_date,omitempty"`
	Status        string    `json:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
	PaymentDate   null.Time `json:"payment_date,omitempty"`
	PaymentMethod string    `json:"payment_method" validate:"omitempty,max=50"`
	Notes         string    `json:"notes"`
}

// Notification DTOs
type CreateNotificationRequest struct {
	UserID  null.Int `json:"user_id" validate:"omitempty,min=1"`
	Title   string   `json:"title" validate:"required,max=100"`
	Message string   `json:"message" validate:"required"`
}

type UpdateNotificationRequest struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Message string `json:"message,omitempty"`
	IsRead  *bool  `json:"is_read,omitempty"`
}

// Query parameters for filtering and pagination
type QueryParams struct {
	Page    int    `query:"page" validate:"omitempty,min=1"`
	Limit   int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search  string `query:"search"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
}

type UserQueryParams struct {
	QueryParams
	RoleID   int   `query:"role_id" validate:"omitempty,min=1"`
	IsActive *bool `query:"is_active"`
}

type StudentQueryParams struct {
	QueryParams
	ClassID  int `query:"class_id" validate:"omitempty,min=1"`
	ParentID int `query:"parent_id" validate:"omitempty,min=1"`
}

type TeacherQueryParams struct {
	QueryParams
	DepartmentID int `query:"department_id" validate:"omitempty,min=1"`
}

type AttendanceQueryParams struct {
	QueryParams
	StudentID  int       `query:"student_id" validate:"omitempty,min=1"`
	ScheduleID int       `query:"schedule_id" validate:"omitempty,min=1"`
	DateFrom   time.Time `query:"date_from"`
	DateTo     time.Time `query:"date_to"`
	Status     string    `query:"status" validate:"omitempty,oneof=Present Absent Late Excused"`
}

type FeeQueryParams struct {
	QueryParams
	StudentID      int    `query:"student_id" validate:"omitempty,min=1"`
	FeeTypeID      int    `query:"fee_type_id" validate:"omitempty,min=1"`
	AcademicYearID int    `query:"academic_year_id" validate:"omitempty,min=1"`
	Status         string `query:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
}
