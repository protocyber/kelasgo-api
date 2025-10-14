package dto

import (
	"time"

	"github.com/google/uuid"
)

// Attendance DTOs
type CreateAttendanceRequest struct {
	StudentID      *uuid.UUID `json:"student_id" validate:"omitempty,uuid"`
	ScheduleID     *uuid.UUID `json:"schedule_id" validate:"omitempty,uuid"`
	Status         string     `json:"status" validate:"required,oneof=present absent late excused"`
	AttendanceDate *time.Time `json:"attendance_date,omitempty"`
	Remarks        *string    `json:"remarks,omitempty"`
}

type UpdateAttendanceRequest struct {
	Status         *string    `json:"status" validate:"omitempty,oneof=present absent late excused"`
	AttendanceDate *time.Time `json:"attendance_date,omitempty"`
	Remarks        *string    `json:"remarks,omitempty"`
}

type AttendanceQueryParams struct {
	QueryParams
	StudentID  *uuid.UUID `query:"student_id" validate:"omitempty,uuid"`
	ScheduleID *uuid.UUID `query:"schedule_id" validate:"omitempty,uuid"`
	DateFrom   *time.Time `query:"date_from"`
	DateTo     *time.Time `query:"date_to"`
	Status     *string    `query:"status" validate:"omitempty,oneof=present absent late excused"`
}
