package usecase

import (
	"context"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/google/uuid"
)

type MasterUseCase struct {
	masterRepo repository.MasterRepository
	schedRepo  repository.ScheduleRepository
}

func NewMasterUseCase(mr repository.MasterRepository, sr repository.ScheduleRepository) *MasterUseCase {
	return &MasterUseCase{
		masterRepo: mr,
		schedRepo:  sr,
	}
}

func (uc *MasterUseCase) CreateMaster(ctx context.Context, master *entity.Master) (*entity.Master, error) {
	if err := uc.masterRepo.Create(ctx, master); err != nil {
		return nil, err
	}

	return master, nil
}

func (uc *MasterUseCase) GetMasterByID(ctx context.Context, id uuid.UUID) (*entity.Master, error) {
	master, err := uc.masterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return master, nil
}

func (uc *MasterUseCase) ListMasters(ctx context.Context, offset, limit int) ([]*entity.Master, error) {
	return uc.masterRepo.List(ctx, offset, limit)
}

func (uc *MasterUseCase) UpdateMaster(ctx context.Context, m *entity.Master) error {
	return uc.masterRepo.Update(ctx, m)
}

func (uc *MasterUseCase) DeleteMaster(ctx context.Context, id uuid.UUID) error {
	return uc.masterRepo.Delete(ctx, id)
}
