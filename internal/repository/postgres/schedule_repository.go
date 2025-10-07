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

type ScheduleRepository struct {
	conn *pgxpool.Pool
}

func NewScheduleRepository(conn *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{conn: conn}
}

func (r *ScheduleRepository) Create(ctx context.Context, s *entity.Schedule) error {
	query := `
		INSERT INTO schedules (master_id, name, type, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id`

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	return r.conn.QueryRow(ctx, query, s.MasterID, s.Name, s.Type, s.StartDate, s.EndDate, now).Scan(&s.ID)
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	query := `
		SELECT id, master_id, name, type, start_date, end_date, created_at, updated_at
		FROM schedules
		WHERE id = $1`

	schedule := &entity.Schedule{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.MasterID,
		&schedule.Name,
		&schedule.Type,
		&schedule.StartDate,
		&schedule.EndDate,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apiErrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *ScheduleRepository) GetByMasterID(ctx context.Context, masterID uuid.UUID) ([]*entity.Schedule, error) {
	query := `
		SELECT id, master_id, name, type, start_date, end_date, created_at, updated_at
		FROM schedules
		WHERE master_id = $1
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*entity.Schedule
	for rows.Next() {
		s := &entity.Schedule{}
		if err := rows.Scan(&s.ID, &s.MasterID, &s.Name, &s.Type, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) GetForDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]*entity.Schedule, error) {
	query := `
		SELECT id, master_id, name, type, start_date, end_date, created_at, updated_at
		FROM schedules
		WHERE master_id = $1 AND start_date <= $2 AND (end_date IS NULL OR end_date >= $2)
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query, masterID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*entity.Schedule
	for rows.Next() {
		s := &entity.Schedule{}
		if err := rows.Scan(&s.ID, &s.MasterID, &s.Name, &s.Type, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *entity.Schedule) error {
	query := `
		UPDATE schedules 
		SET master_id = $1, name = $2, type = $3, start_date = $4, end_date = $5, updated_at = $6
		WHERE id = $7`

	result, err := r.conn.Exec(ctx, query,
		schedule.MasterID,
		schedule.Name,
		schedule.Type,
		schedule.StartDate,
		schedule.EndDate,
		time.Now(),
		schedule.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *ScheduleRepository) List(ctx context.Context, offset, limit int) ([]*entity.Schedule, error) {
	query := `
		SELECT id, master_id, name, type, start_date, end_date, created_at, updated_at
		FROM schedules
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*entity.Schedule
	for rows.Next() {
		schedule := &entity.Schedule{}
		err := rows.Scan(
			schedule.ID,
			schedule.MasterID,
			schedule.Name,
			schedule.Type,
			schedule.StartDate,
			schedule.EndDate,
			schedule.CreatedAt,
			schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *ScheduleRepository) AddDay(ctx context.Context, day *entity.ScheduleDay) error {
	query := `
		INSERT INTO schedule_days (schedule_id, weekday, day_index, start_time, end_time, is_day_off, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id`

	now := time.Now()
	day.CreatedAt = now
	day.UpdatedAt = now

	return r.conn.QueryRow(ctx, query, day.ScheduleID, day.Weekday, day.DayIndex, day.StartTime, day.EndTime, day.IsDayOff, now).Scan(&day.ID)
}

func (r *ScheduleRepository) GetDayByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleDay, error) {
	query := `
		SELECT id, schedule_id, weekday, day_index, start_time, end_time, is_day_off, created_at, updated_at
		FROM schedule_days
		WHERE id = $1`

	day := &entity.ScheduleDay{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&day.ID,
		&day.ScheduleID,
		&day.Weekday,
		&day.DayIndex,
		&day.StartTime,
		&day.EndTime,
		&day.IsDayOff,
		&day.CreatedAt,
		&day.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apiErrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return day, nil
}

func (r *ScheduleRepository) GetDaysByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleDay, error) {
	query := `
		SELECT id, schedule_id, weekday, day_index, start_time, end_time, is_day_off, created_at, updated_at
		FROM schedule_days
		WHERE schedule_id = $1
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []*entity.ScheduleDay
	for rows.Next() {
		day := &entity.ScheduleDay{}
		if err := rows.Scan(
			&day.ID,
			&day.ScheduleID,
			&day.Weekday,
			&day.DayIndex,
			&day.StartTime,
			&day.EndTime,
			&day.IsDayOff,
			&day.CreatedAt,
			&day.UpdatedAt,
		); err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, nil
}

func (r *ScheduleRepository) GetDaysByWeekday(ctx context.Context, masterID uuid.UUID, weekday int) ([]*entity.ScheduleDay, error) {
	query := `
		SELECT d.id, d.schedule_id, d.weekday, d.day_index, d.start_time, d.end_time, d.is_day_off, d.created_at, d.updated_at
		FROM schedule_days d
		JOIN schedules w ON d.schedule_id = w.id
		WHERE w.master_id = $1 AND d.weekday = $2
		ORDER BY w.start_date DESC`

	rows, err := r.conn.Query(ctx, query, masterID, weekday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []*entity.ScheduleDay
	for rows.Next() {
		d := &entity.ScheduleDay{}
		if err := rows.Scan(&d.ID, &d.ScheduleID, &d.Weekday, &d.DayIndex, &d.StartTime, &d.EndTime, &d.IsDayOff, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		days = append(days, d)
	}
	return days, nil
}

func (r *ScheduleRepository) GetDaysByDayIndex(ctx context.Context, masterID uuid.UUID, dayIndex int) ([]*entity.ScheduleDay, error) {
	query := `
		SELECT d.id, d.schedule_id, d.weekday, d.day_index, d.start_time, d.end_time, d.is_day_off, d.created_at, d.updated_at
		FROM schedule_days d
		JOIN schedules w ON d.schedule_id = w.id
		WHERE w.master_id = $1 AND d.day_index = $2
		ORDER BY w.start_date DESC`

	rows, err := r.conn.Query(ctx, query, masterID, dayIndex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []*entity.ScheduleDay
	for rows.Next() {
		d := &entity.ScheduleDay{}
		if err := rows.Scan(&d.ID, &d.ScheduleID, &d.Weekday, &d.DayIndex, &d.StartTime, &d.EndTime, &d.IsDayOff, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		days = append(days, d)
	}
	return days, nil
}

func (r *ScheduleRepository) UpdateDay(ctx context.Context, day *entity.ScheduleDay) error {
	query := `
		UPDATE schedule_days 
		SET schedule_id = $1, weekday = $2, day_index = $3, start_time = $4, end_time = $5, is_day_off = $6, updated_at = $7
		WHERE id = $8`

	result, err := r.conn.Exec(ctx, query,
		day.ScheduleID,
		day.Weekday,
		day.DayIndex,
		day.StartTime,
		day.EndTime,
		day.IsDayOff,
		time.Now(),
		day.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *ScheduleRepository) DeleteDay(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_days WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *ScheduleRepository) GetDaysCount(ctx context.Context, scheduleID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM schedule_days
		WHERE schedule_id = $1
		ORDER BY created_at DESC`

	var count int
	err := r.conn.QueryRow(ctx, query, scheduleID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ScheduleRepository) AddSlot(ctx context.Context, slot *entity.ScheduleSlot) error {
	query := `
		INSERT INTO schedule_slots (schedule_id, date, start_time, end_time, is_day_off, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id`

	now := time.Now()
	slot.CreatedAt = now
	slot.UpdatedAt = now

	return r.conn.QueryRow(ctx, query, slot.ScheduleID, slot.Date, slot.StartTime, slot.EndTime, slot.IsDayOff, now).Scan(&slot.ID)
}

func (r *ScheduleRepository) GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error) {
	query := `
		SELECT id, schedule_id, date, start_time, end_time, is_day_off, created_at, updated_at
		FROM schedule_slots
		WHERE id = $1`

	slot := &entity.ScheduleSlot{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&slot.ID,
		&slot.ScheduleID,
		&slot.Date,
		&slot.StartTime,
		&slot.EndTime,
		&slot.IsDayOff,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apiErrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return slot, nil
}

func (r *ScheduleRepository) GetSlotsByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleSlot, error) {
	query := `
		SELECT id, schedule_id, date, start_time, end_time, is_day_off, created_at, updated_at
		FROM schedule_slots
		WHERE schedule_id = $1
		ORDER BY created_at DESC`

	rows, err := r.conn.Query(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*entity.ScheduleSlot
	for rows.Next() {
		slot := &entity.ScheduleSlot{}
		if err := rows.Scan(
			slot.ID,
			slot.ScheduleID,
			slot.Date,
			slot.StartTime,
			slot.EndTime,
			slot.IsDayOff,
			slot.CreatedAt,
			slot.UpdatedAt,
		); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *ScheduleRepository) GetSlotsByDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]*entity.ScheduleSlot, error) {
	query := `
		SELECT s.id, s.schedule_id, s.date, s.start_time, s.end_time, s.is_day_off, s.created_at, s.updated_at
		FROM schedule_slots s
		JOIN schedules w ON s.schedule_id = w.id
		WHERE w.master_id = $1 AND s.date = $2`

	rows, err := r.conn.Query(ctx, query, masterID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*entity.ScheduleSlot
	for rows.Next() {
		slot := &entity.ScheduleSlot{}
		err := rows.Scan(&slot.ID, &slot.ScheduleID, &slot.Date, &slot.StartTime, &slot.EndTime, &slot.IsDayOff, &slot.CreatedAt, &slot.UpdatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *ScheduleRepository) UpdateSlot(ctx context.Context, day *entity.ScheduleSlot) error {
	query := `
		UPDATE schedule_slots 
		SET schedule_id = $1, date = $2, start_time = $3, end_time = $4, is_day_off = $5, updated_at = $6
		WHERE id = $7`

	result, err := r.conn.Exec(ctx, query,
		day.ScheduleID,
		day.Date,
		day.StartTime,
		day.EndTime,
		day.IsDayOff,
		time.Now(),
		day.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}

func (r *ScheduleRepository) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_slots WHERE id = $1`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apiErrors.ErrNotFound
	}

	return nil
}
