package usecase

import (
	"context"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/google/uuid"
)

type ServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceUseCase(repo repository.ServiceRepository) *ServiceUseCase {
	return &ServiceUseCase{serviceRepo: repo}
}

func (uc *ServiceUseCase) CreateService(ctx context.Context, s *entity.Service) (*entity.Service, error) {
	if err := uc.serviceRepo.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *ServiceUseCase) GetServiceByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	return uc.serviceRepo.GetByID(ctx, id)
}

func (uc *ServiceUseCase) ListServicesByMaster(ctx context.Context, masterID uuid.UUID) ([]*entity.Service, error) {
	return uc.serviceRepo.GetByMasterID(ctx, masterID)
}

func (uc *ServiceUseCase) UpdateService(ctx context.Context, s *entity.Service) error {
	return uc.serviceRepo.Update(ctx, s)
}

func (uc *ServiceUseCase) DeleteService(ctx context.Context, id uuid.UUID) error {
	return uc.serviceRepo.Delete(ctx, id)
}
