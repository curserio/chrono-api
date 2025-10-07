package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	apiErrors "github.com/curserio/chrono-api/internal/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	conn *pgxpool.Pool
}

func NewServiceRepository(conn *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{conn: conn}
}

func (r *ServiceRepository) Create(ctx context.Context, service *entity.Service) error {
	query := `
		INSERT INTO services (master_id, name, description, duration, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id`
	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now

	return r.conn.QueryRow(ctx, query,
		service.MasterID,
		service.Name,
		service.Description,
		service.Duration,
		service.Price,
		now,
	).Scan(&service.ID)
}

func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	query := `
		SELECT id, master_id, name, description, duration, price, created_at, updated_at
		FROM services
		WHERE id = $1`

	service := &entity.Service{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.MasterID,
		&service.Name,
		&service.Description,
		&service.Duration,
		&service.Price,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, err
	}
	return service, nil
}

func (r *ServiceRepository) GetByMasterID(ctx context.Context, masterID uuid.UUID) ([]*entity.Service, error) {
	query := `
		SELECT id, master_id, name, description, duration, price, created_at, updated_at
		FROM services
		WHERE master_id = $1
		ORDER BY name`

	rows, err := r.conn.Query(ctx, query, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*entity.Service
	for rows.Next() {
		s := &entity.Service{}
		if err := rows.Scan(
			&s.ID,
			&s.MasterID,
			&s.Name,
			&s.Description,
			&s.Duration,
			&s.Price,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, rows.Err()
}

func (r *ServiceRepository) Update(ctx context.Context, service *entity.Service) error {
	query := `
		UPDATE services
		SET name=$1, description=$2, duration=$3, price=$4, updated_at=$5
		WHERE id=$6`

	result, err := r.conn.Exec(ctx, query,
		service.Name,
		service.Description,
		service.Duration,
		service.Price,
		time.Now(),
		service.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}

func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM services WHERE id=$1`
	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}
