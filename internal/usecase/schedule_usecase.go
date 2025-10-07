package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/dto"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/curserio/chrono-api/pkg/timeutil"
	"github.com/google/uuid"
)

type ScheduleUseCase struct {
	repo repository.ScheduleRepository
}

func NewScheduleUseCase(repo repository.ScheduleRepository) *ScheduleUseCase {
	return &ScheduleUseCase{repo: repo}
}

func (uc *ScheduleUseCase) CreateSchedule(ctx context.Context, schedule *entity.Schedule, days []*entity.ScheduleDay) (*entity.Schedule, error) {
	if err := uc.repo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	// Добавляем дни недели
	for _, d := range days {
		d.ScheduleID = schedule.ID
		if err := uc.repo.AddDay(ctx, d); err != nil {
			return nil, fmt.Errorf("add schedule day: %w", err)
		}
	}

	return schedule, nil
}

func (uc *ScheduleUseCase) GetScheduleByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	schedule, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (uc *ScheduleUseCase) GetSchedulesByMaster(ctx context.Context, masterID uuid.UUID) ([]*entity.Schedule, error) {
	schedules, err := uc.repo.GetByMasterID(ctx, masterID)
	if err != nil {
		return nil, fmt.Errorf("get schedules by master: %w", err)
	}

	return schedules, nil
}

func (uc *ScheduleUseCase) ListSchedules(ctx context.Context, offset, limit int) ([]*entity.Schedule, error) {
	return uc.repo.List(ctx, offset, limit)
}

func (uc *ScheduleUseCase) UpdateSchedule(ctx context.Context, s *entity.Schedule) error {
	if err := uc.repo.Update(ctx, s); err != nil {
		return fmt.Errorf("update schedule: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) AddDay(ctx context.Context, day *entity.ScheduleDay) error {
	if err := uc.repo.AddDay(ctx, day); err != nil {
		return fmt.Errorf("add day: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) GetDayByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleDay, error) {
	day, err := uc.repo.GetDayByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get day: %w", err)
	}
	return day, nil
}

func (uc *ScheduleUseCase) GetDaysBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleDay, error) {
	days, err := uc.repo.GetDaysByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("get days by schedule: %w", err)
	}
	return days, nil
}

func (uc *ScheduleUseCase) GetDaysByWeekday(ctx context.Context, masterID uuid.UUID, weekday int) ([]*entity.ScheduleDay, error) {
	days, err := uc.repo.GetDaysByWeekday(ctx, masterID, weekday)
	if err != nil {
		return nil, fmt.Errorf("get days by weekday: %w", err)
	}
	return days, nil
}

func (uc *ScheduleUseCase) UpdateDay(ctx context.Context, day *entity.ScheduleDay) error {
	if err := uc.repo.UpdateDay(ctx, day); err != nil {
		return fmt.Errorf("update day: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) DeleteDay(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.DeleteDay(ctx, id); err != nil {
		return fmt.Errorf("delete day: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) AddSlot(ctx context.Context, slot *entity.ScheduleSlot) error {
	if err := uc.repo.AddSlot(ctx, slot); err != nil {
		return fmt.Errorf("add slot: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error) {
	slot, err := uc.repo.GetSlotByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get slot: %w", err)
	}
	return slot, nil
}

func (uc *ScheduleUseCase) GetSlotsBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]*entity.ScheduleSlot, error) {
	slots, err := uc.repo.GetSlotsByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("get slots by schedule: %w", err)
	}
	return slots, nil
}

func (uc *ScheduleUseCase) GetSlotsByDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]*entity.ScheduleSlot, error) {
	// Нормализуем дату, чтобы отсеять время
	date = timeutil.NormalizeDate(date)

	slots, err := uc.repo.GetSlotsByDate(ctx, masterID, date)
	if err != nil {
		return nil, fmt.Errorf("get slots by date: %w", err)
	}
	return slots, nil
}

func (uc *ScheduleUseCase) UpdateSlot(ctx context.Context, slot *entity.ScheduleSlot) error {
	if err := uc.repo.UpdateSlot(ctx, slot); err != nil {
		return fmt.Errorf("update slot: %w", err)
	}
	return nil
}

func (uc *ScheduleUseCase) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.DeleteSlot(ctx, id); err != nil {
		return fmt.Errorf("delete slot: %w", err)
	}
	return nil
}

// GetScheduleForDate returns the effective schedule for a given master and date.
//
// The logic is as follows:
// 1. Check for override slots (specific date overrides). If present, return them immediately (source = "override").
// 2. Otherwise, determine which base schedule applies to the given date (weekly or cyclic).
// 3. Based on schedule type:
//   - Weekly: use the weekday to select schedule_days.
//   - Cyclic: calculate day_index from (date - start_date) % cycle_length.
//
// 4. Return the resulting working hours or mark as a day off.
//
// The response is always a slice of ScheduleForDateResponse objects —
// even if there are no active slots for that date.
func (uc *ScheduleUseCase) GetScheduleForDate(ctx context.Context, masterID uuid.UUID, date time.Time) ([]dto.ScheduleForDateResponse, error) {
	// Normalize date to midnight UTC to ensure consistent lookups.
	date = timeutil.NormalizeDate(date)

	// Check for override slots for the exact date.
	overrideSlots, err := uc.repo.GetSlotsByDate(ctx, masterID, date)
	if err != nil {
		return nil, fmt.Errorf("get slots by date: %w", err)
	}
	if len(overrideSlots) > 0 {
		out := make([]dto.ScheduleForDateResponse, 0, len(overrideSlots))
		for _, s := range overrideSlots {
			var st, et *string
			if !s.IsDayOff {
				if s.StartTime != nil {
					tmp := formatTimeOfDay(*s.StartTime)
					st = &tmp
				}
				if s.EndTime != nil {
					tmp := formatTimeOfDay(*s.EndTime)
					et = &tmp
				}
			}
			out = append(out, dto.ScheduleForDateResponse{
				MasterID:  masterID,
				Date:      date.Format(time.DateOnly),
				StartTime: st,
				EndTime:   et,
				IsDayOff:  s.IsDayOff,
				Source:    "override",
			})
		}
		return out, nil
	}

	// Find the active schedule for the given date (may return multiple, pick the latest).
	schedules, err := uc.repo.GetForDate(ctx, masterID, date)
	if err != nil {
		return nil, fmt.Errorf("get schedules: %w", err)
	}
	if len(schedules) == 0 {
		return nil, nil // no schedule at all
	}
	schedule := schedules[0]

	// Determine which days to fetch depending on schedule type.
	var days []*entity.ScheduleDay
	switch schedule.Type {
	case entity.ScheduleTypeCyclic:
		// Calculate which cycle day corresponds to the current date.
		cycleLength, err := uc.repo.GetDaysCount(ctx, schedule.ID)
		if err != nil {
			return nil, fmt.Errorf("get schedule cycle length: %w", err)
		}

		daysSinceStart := int(date.Truncate(24*time.Hour).Sub(schedule.StartDate.Truncate(24*time.Hour)).Hours() / 24)
		if cycleLength <= 0 {
			return nil, fmt.Errorf("invalid cycle length for schedule %s", schedule.ID)
		}
		dayIndex := (daysSinceStart % cycleLength) + 1

		days, err = uc.repo.GetDaysByDayIndex(ctx, masterID, dayIndex)
		if err != nil {
			return nil, fmt.Errorf("get schedule days by day index: %w", err)
		}

	case entity.ScheduleTypeWeekly:
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7 // make Sunday = 7 to match DB convention
		}
		days, err = uc.repo.GetDaysByWeekday(ctx, masterID, weekday)
		if err != nil {
			return nil, fmt.Errorf("get schedule days by weekday: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported schedule type: %s", schedule.Type)
	}

	if len(days) == 0 {
		// no working hours found — return as day off
		return []dto.ScheduleForDateResponse{{
			MasterID: masterID,
			Date:     date.Format("2006-01-02"),
			IsDayOff: true,
			Source:   schedule.Name,
		}}, nil
	}

	// Build response DTOs
	out := make([]dto.ScheduleForDateResponse, 0, len(days))
	for _, d := range days {
		res := dto.ScheduleForDateResponse{
			MasterID: masterID,
			Date:     date.Format("2006-01-02"),
			IsDayOff: d.IsDayOff,
			Source:   schedule.Name,
		}
		if d.StartTime != nil && d.EndTime != nil {
			st := formatTimeOfDay(*d.StartTime)
			et := formatTimeOfDay(*d.EndTime)
			res.StartTime = &st
			res.EndTime = &et
		}
		out = append(out, res)
	}
	return out, nil
}

// GetScheduleForRange returns the master's schedule for a given date range (inclusive).
// If no data is found for a particular date, a placeholder entry is returned with IsDayOff = true and Source = "none".
// This simplifies rendering on the frontend side.
func (uc *ScheduleUseCase) GetScheduleForRange(ctx context.Context, masterID uuid.UUID, fromDate, toDate time.Time) ([]dto.ScheduleForDateResponse, error) {
	results := make([]dto.ScheduleForDateResponse, 0)
	current := fromDate

	for !current.After(toDate) {
		daily, err := uc.GetScheduleForDate(ctx, masterID, current)
		if err != nil {
			return nil, fmt.Errorf("get schedule for date %s: %w", current.Format(time.DateOnly), err)
		}

		if len(daily) == 0 {
			results = append(results, dto.ScheduleForDateResponse{
				MasterID: masterID,
				Date:     current.Format(time.DateOnly),
				IsDayOff: true,
				Source:   "none",
			})
		} else {
			results = append(results, daily...)
		}

		current = current.AddDate(0, 0, 1)
	}

	return results, nil
}

// formatTimeOfDay returns "HH:mm" representation of a time-of-day field.
func formatTimeOfDay(t time.Time) string {
	return t.Format("15:04")
}
