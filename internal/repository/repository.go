package repository

import (
	"context"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/mock_repository.go -package=mock github.com/curserio/chrono-api/internal/repository MasterRepository,ScheduleRepository,ServiceRepository,BookingRepository,ClientRepository

type MasterRepository interface {
	Create(ctx context.Context, master *entity.Master) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Master, error)
	Update(ctx context.Context, master *entity.Master) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*entity.Master, error)
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	GetByMasterID(ctx context.Context, masterID uuid.UUID) ([]*entity.Schedule, error)
	GetForDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]*entity.Schedule, error)
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*entity.Schedule, error)

	AddDay(ctx context.Context, day *entity.ScheduleDay) error
	GetDayByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleDay, error)
	GetDaysByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleDay, error)
	GetDaysByWeekday(ctx context.Context, masterID uuid.UUID, weekday int) ([]*entity.ScheduleDay, error)
	GetDaysByDayIndex(ctx context.Context, masterID uuid.UUID, dayIndex int) ([]*entity.ScheduleDay, error)
	UpdateDay(ctx context.Context, day *entity.ScheduleDay) error
	DeleteDay(ctx context.Context, id uuid.UUID) error

	GetDaysCount(ctx context.Context, scheduleID uuid.UUID) (int, error)

	AddSlot(ctx context.Context, slot *entity.ScheduleSlot) error
	GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error)
	GetSlotsByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleSlot, error)
	GetSlotsByDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]*entity.ScheduleSlot, error)
	UpdateSlot(ctx context.Context, day *entity.ScheduleSlot) error
	DeleteSlot(ctx context.Context, id uuid.UUID) error
}

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetByMasterID(ctx context.Context, masterID uuid.UUID) ([]*entity.Service, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetByMasterID(ctx context.Context, masterID uuid.UUID, from, to time.Time) ([]*entity.Booking, error)
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*entity.Booking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BookingStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ClientRepository interface {
	Create(ctx context.Context, client *entity.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Client, error)
	Update(ctx context.Context, client *entity.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*entity.Client, error)
}
