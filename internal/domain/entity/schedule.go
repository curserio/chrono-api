package entity

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleType string

const (
	ScheduleTypeWeekly ScheduleType = "weekly"
	ScheduleTypeCyclic ScheduleType = "cyclic"
	ScheduleTypeCustom ScheduleType = "custom"
)

type Schedule struct {
	ID        uuid.UUID    `json:"id"`
	MasterID  uuid.UUID    `json:"master_id"`
	Name      string       `json:"name"`
	Type      ScheduleType `json:"type"`
	StartDate time.Time    `json:"start_date"`
	EndDate   *time.Time   `json:"end_date"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type ScheduleDay struct {
	ID         uuid.UUID  `json:"id"`
	ScheduleID uuid.UUID  `json:"schedule_id"`
	Weekday    *int       `json:"weekday,omitempty"`
	DayIndex   *int       `json:"day_index,omitempty"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	IsDayOff   bool       `json:"is_day_off"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ScheduleSlot struct {
	ID         uuid.UUID  `json:"id"`
	ScheduleID uuid.UUID  `json:"schedule_id"`
	Date       time.Time  `json:"date"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	IsDayOff   bool       `json:"is_day_off"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
