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

type MasterRepository struct {
	conn *pgxpool.Pool
}

func NewMasterRepository(conn *pgxpool.Pool) *MasterRepository {
	return &MasterRepository{conn: conn}
}

func (r *MasterRepository) Create(ctx context.Context, master *entity.Master) error {
	query := `
		INSERT INTO masters (name, email, phone, telegram_id, telegram_username, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id`

	now := time.Now()
	master.CreatedAt = now
	master.UpdatedAt = now

	err := r.conn.QueryRow(ctx, query,
		master.Name,
		master.Email,
		master.Phone,
		master.TelegramID,
		master.TelegramUsername,
		master.Description,
		now,
	).Scan(&master.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *MasterRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Master, error) {
	query := `
		SELECT id, name, email, phone, telegram_id, telegram_username, description, created_at, updated_at
		FROM masters
		WHERE id = $1`

	master := &entity.Master{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&master.ID,
		&master.Name,
		&master.Email,
		&master.Phone,
		&master.TelegramID,
		&master.TelegramUsername,
		&master.Description,
		&master.CreatedAt,
		&master.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apiErrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return master, nil
}

func (r *MasterRepository) Update(ctx context.Context, master *entity.Master) error {
	query := `
		UPDATE masters 
		SET name = $1, email = $2, phone = $3, telegram_id = $4, telegram_username = $5, description = $6, updated_at = $7
		WHERE id = $8`

	result, err := r.conn.Exec(ctx, query,
		master.Name,
		master.Email,
		master.Phone,
		master.TelegramID,
		master.TelegramUsername,
		master.Description,
		time.Now(),
		master.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *MasterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM masters WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *MasterRepository) List(ctx context.Context, offset, limit int) ([]*entity.Master, error) {
	query := `
		SELECT id, name, email, phone, telegram_id, telegram_username, description, created_at, updated_at
		FROM masters
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var masters []*entity.Master
	for rows.Next() {
		master := &entity.Master{}
		err := rows.Scan(
			&master.ID,
			&master.Name,
			&master.Email,
			&master.Phone,
			&master.TelegramID,
			&master.TelegramUsername,
			&master.Description,
			&master.CreatedAt,
			&master.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return masters, nil
}
