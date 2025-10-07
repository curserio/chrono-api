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

type ClientRepository struct {
	conn *pgxpool.Pool
}

func NewClientRepository(conn *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{conn: conn}
}

func (r *ClientRepository) Create(ctx context.Context, client *entity.Client) error {
	query := `
		INSERT INTO clients (name, email, phone, telegram_id, telegram_username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id`

	now := time.Now()
	client.CreatedAt = now
	client.UpdatedAt = now

	err := r.conn.QueryRow(ctx, query,
		client.Name,
		client.Email,
		client.Phone,
		client.TelegramID,
		client.TelegramUsername,
		now,
	).Scan(&client.ID)

	return err
}

func (r *ClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	query := `
		SELECT id, name, email, phone, telegram_id, telegram_username, created_at, updated_at
		FROM clients
		WHERE id = $1`

	client := &entity.Client{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.Phone,
		&client.TelegramID,
		&client.TelegramUsername,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, err
	}
	return client, nil
}

func (r *ClientRepository) Update(ctx context.Context, client *entity.Client) error {
	query := `
		UPDATE clients
		SET name = $1, email = $2, phone = $3, telegram_id = $4, telegram_username = $5, updated_at = $6
		WHERE id = $7`

	result, err := r.conn.Exec(ctx, query,
		client.Name,
		client.Email,
		client.Phone,
		client.TelegramID,
		client.TelegramUsername,
		time.Now(),
		client.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}

func (r *ClientRepository) List(ctx context.Context, offset, limit int) ([]*entity.Client, error) {
	query := `
		SELECT id, name, email, phone, telegram_id, telegram_username, created_at, updated_at
		FROM clients
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]*entity.Client, 0)
	for rows.Next() {
		client := &entity.Client{}
		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Email,
			&client.Phone,
			&client.TelegramID,
			&client.TelegramUsername,
			&client.CreatedAt,
			&client.UpdatedAt,
		); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, rows.Err()
}
