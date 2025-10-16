package dto

import (
	"github.com/google/uuid"
)

// Schedule DTOs
type CreateScheduleRequest struct {
	ClassSubjectID *uuid.UUID `json:"class_subject_id" validate:"omitempty,uuid"`
	DayOfWeek      DayOfWeek  `json:"day_of_week" validate:"required,oneof=senin selasa rabu kamis jumat sabtu minggu"`
	StartTime      string     `json:"start_time" validate:"required"`
	EndTime        string     `json:"end_time" validate:"required"`
	Room           *string    `json:"room" validate:"omitempty,max=50"`
}

type UpdateScheduleRequest struct {
	ClassSubjectID *uuid.UUID `json:"class_subject_id" validate:"omitempty,uuid"`
	DayOfWeek      *DayOfWeek `json:"day_of_week" validate:"omitempty,oneof=senin selasa rabu kamis jumat sabtu minggu"`
	StartTime      *string    `json:"start_time,omitempty"`
	EndTime        *string    `json:"end_time,omitempty"`
	Room           *string    `json:"room" validate:"omitempty,max=50"`
}
