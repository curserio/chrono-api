package usecase

import (
	"context"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/google/uuid"
)

type BookingUseCase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(repo repository.BookingRepository) *BookingUseCase {
	return &BookingUseCase{bookingRepo: repo}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, booking *entity.Booking) (*entity.Booking, error) {
	if err := uc.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (uc *BookingUseCase) GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	return uc.bookingRepo.GetByID(ctx, id)
}

func (uc *BookingUseCase) GetBookingsByMaster(ctx context.Context, masterID uuid.UUID, from, to time.Time) ([]*entity.Booking, error) {
	return uc.bookingRepo.GetByMasterID(ctx, masterID, from, to)
}

func (uc *BookingUseCase) GetBookingsByClient(ctx context.Context, clientID uuid.UUID) ([]*entity.Booking, error) {
	return uc.bookingRepo.GetByClientID(ctx, clientID)
}

func (uc *BookingUseCase) UpdateBookingStatus(ctx context.Context, id uuid.UUID, status entity.BookingStatus) error {
	return uc.bookingRepo.UpdateStatus(ctx, id, status)
}

func (uc *BookingUseCase) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	return uc.bookingRepo.Delete(ctx, id)
}
