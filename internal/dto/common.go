package dto

// Enums to match database schema
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type DayOfWeek string

const (
	DayMonday    DayOfWeek = "senin"
	DayTuesday   DayOfWeek = "selasa"
	DayWednesday DayOfWeek = "rabu"
	DayThursday  DayOfWeek = "kamis"
	DayFriday    DayOfWeek = "jumat"
	DaySaturday  DayOfWeek = "sabtu"
	DaySunday    DayOfWeek = "minggu"
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

// Query parameters for filtering and pagination
type QueryParams struct {
	Page    int    `query:"page" validate:"omitempty,min=1"`
	Limit   int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search  string `query:"search"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
}
