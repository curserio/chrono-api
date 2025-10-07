package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/dto"
	"github.com/curserio/chrono-api/internal/errors"
	"github.com/curserio/chrono-api/internal/repository/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMasterUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockMasterRepository(ctrl)
	useCase := NewMasterUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New()
		expectedMaster := &entity.Master{
			ID:        id,
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.EXPECT().
			GetByID(ctx, id).
			Return(expectedMaster, nil)

		master, err := useCase.GetMasterByID(ctx, expectedMaster.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedMaster, master)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		id := uuid.New()
		mockRepo.EXPECT().
			GetByID(ctx, id).
			Return(nil, errors.ErrNotFound)

		master, err := useCase.GetMasterByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, err)
		assert.Nil(t, master)
	})
}

func TestMasterUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockMasterRepository(ctrl)
	useCase := NewMasterUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		req := &dto.CreateMasterRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Phone:    "+1234567890",
			Password: "123456",
		}

		createdMaster := &entity.Master{
			ID:    uuid.New(),
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, master *entity.Master) error {
				master.ID = createdMaster.ID
				return nil
			})

		master, err := useCase.CreateMaster(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, createdMaster, master)
	})
}
