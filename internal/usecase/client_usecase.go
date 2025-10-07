package usecase

import (
	"context"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/google/uuid"
)

type ClientUseCase struct {
	clientRepo repository.ClientRepository
}

func NewClientUseCase(repo repository.ClientRepository) *ClientUseCase {
	return &ClientUseCase{clientRepo: repo}
}

func (uc *ClientUseCase) CreateClient(ctx context.Context, c *entity.Client) (*entity.Client, error) {
	if err := uc.clientRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *ClientUseCase) GetClientByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	return uc.clientRepo.GetByID(ctx, id)
}

func (uc *ClientUseCase) ListClients(ctx context.Context, offset, limit int) ([]*entity.Client, error) {
	return uc.clientRepo.List(ctx, offset, limit)
}

func (uc *ClientUseCase) UpdateClient(ctx context.Context, c *entity.Client) error {
	return uc.clientRepo.Update(ctx, c)
}

func (uc *ClientUseCase) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return uc.clientRepo.Delete(ctx, id)
}
