package model

import (
	"github.com/google/uuid"
)

// Schedule represents the schedules table
type Schedule struct {
	BaseModel
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ClassSubjectID *uuid.UUID `gorm:"type:uuid;index" json:"class_subject_id,omitempty"`
	DayOfWeek      string     `gorm:"size:15;check:day_of_week IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday')" json:"day_of_week"`
	StartTime      string     `gorm:"type:time" json:"start_time"`
	EndTime        string     `gorm:"type:time" json:"end_time"`
	Room           *string    `gorm:"size:50" json:"room,omitempty"`

	// Relationships
	ClassSubject *ClassSubject `gorm:"foreignKey:ClassSubjectID;constraint:OnDelete:CASCADE" json:"class_subject,omitempty"`
	Attendance   []Attendance  `gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE" json:"attendance,omitempty"`
}

// TableName returns the table name for Schedule
func (Schedule) TableName() string {
	return "schedules"
}
