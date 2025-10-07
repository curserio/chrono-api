package dto

import (
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/errors"
	"github.com/google/uuid"
)

type CreateMasterRequest struct {
	Name             string  `json:"name" validate:"required,min=2,max=255"`
	Email            *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone            *string `json:"phone,omitempty" validate:"omitempty,e164"`
	TelegramID       *int64  `json:"telegram_id,omitempty" validate:"omitempty,gte=1"`
	TelegramUsername *string `json:"telegram_username,omitempty" validate:"omitempty,min=1,max=255"`
	Description      *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	City             *string `json:"city" validate:"omitempty,min=1,max=100"`
	Timezone         string  `json:"timezone" validate:"omitempty,timezone"`
	Language         string  `json:"language" validate:"omitempty,len=2"`
}

type UpdateMasterRequest struct {
	Name             string  `json:"name" validate:"omitempty,min=2,max=255"`
	Email            *string `json:"email" validate:"omitempty,email"`
	Phone            *string `json:"phone" validate:"omitempty,e164"`
	TelegramID       *int64  `json:"telegram_id" validate:"omitempty,gte=1"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty,min=1,max=255"`
	Description      *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	City             *string `json:"city" validate:"omitempty,min=1,max=100"`
	Timezone         string  `json:"timezone" validate:"omitempty,timezone"`
	Language         string  `json:"language" validate:"omitempty,len=2"`
}

type CreateClientRequest struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`
	Email            *string `json:"email" validate:"omitempty,email"`
	Phone            *string `json:"phone" validate:"omitempty,e164"`
	TelegramID       *int64  `json:"telegram_id" validate:"omitempty,gte=1"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty,min=1,max=255"`
	City             *string `json:"city" validate:"omitempty,min=1,max=100"`
	Timezone         string  `json:"timezone" validate:"omitempty,timezone"` // e.g. Europe/Moscow
	Language         string  `json:"language" validate:"omitempty,len=2"`    // ISO 639-1
}

type UpdateClientRequest struct {
	Name             string  `json:"name" validate:"omitempty,min=2,max=100"`
	Email            *string `json:"email" validate:"omitempty,email"`
	Phone            *string `json:"phone" validate:"omitempty,e164"`
	TelegramID       *int64  `json:"telegram_id" validate:"omitempty,gte=1"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty,min=1,max=255"`
	City             *string `json:"city" validate:"omitempty,min=1,max=100"`
	Timezone         string  `json:"timezone" validate:"omitempty,timezone"`
	Language         string  `json:"language" validate:"omitempty,len=2"`
}

type CreateServiceRequest struct {
	MasterID    uuid.UUID `json:"master_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description" validate:"max=500"`
	Duration    int       `json:"duration" validate:"required,min=1,max=1440"` // макс. 24 часа
	Price       float64   `json:"price" validate:"required,min=0"`
}

type UpdateServiceRequest struct {
	MasterID    uuid.UUID `json:"master_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Description string    `json:"description" validate:"max=500"`
	Duration    int       `json:"duration" validate:"required,min=1,max=1440"` // макс. 24 часа
	Price       float64   `json:"price" validate:"required,min=0"`
}

type CreateBookingRequest struct {
	MasterID  uuid.UUID `json:"master_id" validate:"required"`
	ServiceID uuid.UUID `json:"service_id" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	StartTime string    `json:"start_time" validate:"required"`
	EndTime   string    `json:"end_time" validate:"required"`
	ClientID  uuid.UUID `json:"client_id" validate:"required"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

func (b BookingStatus) ToEntity() (entity.BookingStatus, error) {
	switch b {
	case BookingStatusPending:
		return entity.BookingStatusPending, nil
	case BookingStatusConfirmed:
		return entity.BookingStatusConfirmed, nil
	case BookingStatusCompleted:
		return entity.BookingStatusCompleted, nil
	case BookingStatusCancelled:
		return entity.BookingStatusCancelled, nil
	default:
		return entity.BookingStatusPending, errors.ErrBookingStatusInvalid
	}
}

type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" validate:"required,oneof=pending confirmed completed cancelled"`
}

type ListBookingsRequest struct {
	MasterID uuid.UUID `query:"master_id" validate:"required"`
	FromDate time.Time `query:"from_date" validate:"required"`
	ToDate   time.Time `query:"to_date" validate:"required,gtfield=FromDate"` // должно быть после FromDate
}

type ScheduleType string

const (
	ScheduleTypeWeekly ScheduleType = "weekly"
	ScheduleTypeCyclic ScheduleType = "cyclic"
	ScheduleTypeCustom ScheduleType = "custom"
)

func (b ScheduleType) ToEntity() (entity.ScheduleType, error) {
	switch b {
	case ScheduleTypeWeekly:
		return entity.ScheduleTypeWeekly, nil
	case ScheduleTypeCyclic:
		return entity.ScheduleTypeCyclic, nil
	case ScheduleTypeCustom:
		return entity.ScheduleTypeCustom, nil
	default:
		return entity.ScheduleTypeWeekly, errors.ErrScheduleTypeInvalid
	}
}

type CreateScheduleRequest struct {
	MasterID  uuid.UUID     `json:"master_id" validate:"required"`
	Name      string        `json:"name" validate:"required,min=1,max=100"`
	Type      ScheduleType  `json:"type" validate:"required,oneof=weekly cyclic custom"`
	StartDate *time.Time    `json:"start_date,omitempty"`
	EndDate   *time.Time    `json:"end_date,omitempty"`
	Days      []ScheduleDay `json:"days" validate:"required,dive"`
}

type ScheduleDay struct {
	Weekday   *int    `json:"weekday,omitempty" validate:"omitempty,min=1,max=7"` // 1 — понедельник, 7 — воскресенье
	DayIndex  *int    `json:"day_index,omitempty" validate:"omitempty,min=1,max=31"`
	StartTime *string `json:"start_time" validate:"omitempty"` // формат "15:04"
	EndTime   *string `json:"end_time" validate:"omitempty"`
	IsDayOff  bool    `json:"is_day_off"`
}

type UpdateScheduleRequest struct {
	MasterID  uuid.UUID     `json:"master_id" validate:"required"`
	Name      string        `json:"name" validate:"required,min=1,max=100"`
	Type      ScheduleType  `json:"type" validate:"required,oneof=weekly cyclic custom"`
	StartDate *time.Time    `json:"start_date,omitempty"`
	EndDate   *time.Time    `json:"end_date,omitempty"`
	Days      []ScheduleDay `json:"days" validate:"required,dive"`
}

type CreateScheduleDayRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" validate:"required"`
	Weekday    *int      `json:"weekday,omitempty" validate:"omitempty,min=1,max=7"` // 1 — понедельник, 7 — воскресенье
	DayIndex   *int      `json:"day_index,omitempty" validate:"omitempty,min=1,max=31"`
	StartTime  *string   `json:"start_time" validate:"omitempty"` // формат "15:04"
	EndTime    *string   `json:"end_time" validate:"omitempty"`
	IsDayOff   bool      `json:"is_day_off"`
}

type UpdateScheduleDayRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" validate:"required"`
	Weekday    *int      `json:"weekday,omitempty" validate:"omitempty,min=1,max=7"` // 1 — понедельник, 7 — воскресенье
	DayIndex   *int      `json:"day_index,omitempty" validate:"omitempty,min=1,max=31"`
	StartTime  *string   `json:"start_time" validate:"omitempty"` // формат "15:04"
	EndTime    *string   `json:"end_time" validate:"omitempty"`
	IsDayOff   bool      `json:"is_day_off"`
}

type CreateScheduleSlotRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
	StartTime  *string   `json:"start_time,omitempty" validate:"omitempty,required"`
	EndTime    *string   `json:"end_time,omitempty" validate:"omitempty,required"`
	IsDayOff   bool      `json:"is_day_off"`
}

type UpdateScheduleSlotRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
	StartTime  *string   `json:"start_time,omitempty" validate:"omitempty,required"`
	EndTime    *string   `json:"end_time,omitempty" validate:"omitempty,required"`
	IsDayOff   bool      `json:"is_day_off"`
}

// ScheduleForDateResponse — финальный ответ, если запрашивается расписание на конкретную дату.
// Если есть override — используется он, иначе возвращается имя расписания.
type ScheduleForDateResponse struct {
	MasterID  uuid.UUID `json:"master_id"`
	Date      string    `json:"date"`
	StartTime *string   `json:"start_time,omitempty"`
	EndTime   *string   `json:"end_time,omitempty"`
	IsDayOff  bool      `json:"is_day_off"`
	Source    string    `json:"source"` // "override" или имя расписания
}
