package postgres

import (
	"context"
	"errors"
	"time"

	apiErrors "github.com/curserio/chrono-api/internal/errors"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	conn *pgxpool.Pool
}

func NewBookingRepository(conn *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{conn: conn}
}

func (r *BookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	query := `
		INSERT INTO bookings (master_id, client_id, service_id, start_time, end_time, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$7)
		RETURNING id`

	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	return r.conn.QueryRow(ctx, query,
		booking.MasterID,
		booking.ClientID,
		booking.ServiceID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		now,
	).Scan(&booking.ID)
}

func (r *BookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	query := `
		SELECT id, master_id, client_id, service_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE id = $1`

	b := &entity.Booking{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.MasterID,
		&b.ClientID,
		&b.ServiceID,
		&b.StartTime,
		&b.EndTime,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *BookingRepository) GetByMasterID(ctx context.Context, masterID uuid.UUID, from, to time.Time) ([]*entity.Booking, error) {
	query := `
		SELECT id, master_id, client_id, service_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE master_id=$1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time`

	rows, err := r.conn.Query(ctx, query, masterID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		b := &entity.Booking{}
		if err := rows.Scan(
			&b.ID,
			&b.MasterID,
			&b.ClientID,
			&b.ServiceID,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	return bookings, rows.Err()
}

func (r *BookingRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*entity.Booking, error) {
	query := `
		SELECT id, master_id, client_id, service_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE client_id=$1
		ORDER BY start_time`

	rows, err := r.conn.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		b := &entity.Booking{}
		if err := rows.Scan(
			&b.ID,
			&b.MasterID,
			&b.ClientID,
			&b.ServiceID,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	return bookings, rows.Err()
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BookingStatus) error {
	query := `
		UPDATE bookings
		SET status=$1, updated_at=$2
		WHERE id=$3`

	result, err := r.conn.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}

func (r *BookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM bookings WHERE id=$1`
	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}
	return nil
}
