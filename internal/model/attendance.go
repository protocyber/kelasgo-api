package model

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceStatus represents the attendance status enum
type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "Present"
	AttendanceAbsent  AttendanceStatus = "Absent"
	AttendanceLate    AttendanceStatus = "Late"
	AttendanceExcused AttendanceStatus = "Excused"
)

// Attendance represents the attendance table
type Attendance struct {
	BaseModel
	TenantID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID      *uuid.UUID       `gorm:"type:uuid;index" json:"student_id,omitempty"`
	ScheduleID     *uuid.UUID       `gorm:"type:uuid;index" json:"schedule_id,omitempty"`
	Status         AttendanceStatus `gorm:"type:attendance_status_enum;default:'Present'" json:"status"`
	AttendanceDate time.Time        `gorm:"type:date;default:CURRENT_DATE" json:"attendance_date"`
	Remarks        *string          `gorm:"type:text" json:"remarks,omitempty"`

	// Relationships
	Student  *Student  `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	Schedule *Schedule `gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE" json:"schedule,omitempty"`
}

// TableName returns the table name for Attendance
func (Attendance) TableName() string {
	return "attendance"
}
