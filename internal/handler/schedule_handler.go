package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/dto"
	"github.com/curserio/chrono-api/internal/errors"
	"github.com/curserio/chrono-api/internal/infrastructure/http/server"
	"github.com/curserio/chrono-api/internal/usecase"
	"github.com/curserio/chrono-api/pkg/logger"
	"github.com/curserio/chrono-api/pkg/timeutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ScheduleHandler struct {
	scheduleUseCase *usecase.ScheduleUseCase
}

func NewScheduleHandler(s *server.Server, uc *usecase.ScheduleUseCase) {
	handler := &ScheduleHandler{scheduleUseCase: uc}

	group := s.NewGroup("/api/v1/schedules")
	group.POST("", handler.CreateSchedule)
	group.GET("", handler.ListSchedules)
	group.GET("/:id", handler.GetSchedule)
	group.PUT("/:id", handler.UpdateSchedule)
	group.DELETE("/:id", handler.DeleteSchedule)

	group.GET("/master/:master_id/date/:date", handler.GetScheduleForDate)
	group.GET("/master/:master_id/range", handler.GetScheduleForRange)

	group.POST("/:id/days", handler.AddDay)
	group.GET("/:id/days", handler.ListDays)
	group.PUT("/days/:id", handler.UpdateDay)
	group.DELETE("/days/:id", handler.DeleteDay)

	group.POST("/:id/slots", handler.AddSlot)
	group.GET("/:id/slots", handler.ListSlots)
	group.PUT("/slots/:id", handler.UpdateSlot)
	group.DELETE("/slots/:id", handler.DeleteSlot)
}

func (h *ScheduleHandler) CreateSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	req := new(dto.CreateScheduleRequest)
	if err := c.Bind(req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid body", err)
	}

	if err := c.Validate(req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Validation failed", err)
	}

	scheduleType, err := req.Type.ToEntity()
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid type value", err)
	}

	schedule := &entity.Schedule{
		MasterID:  req.MasterID,
		Name:      req.Name,
		Type:      scheduleType,
		StartDate: time.Now(),
		EndDate:   req.EndDate,
	}
	if req.StartDate != nil {
		schedule.StartDate = *req.StartDate
	}

	days := make([]*entity.ScheduleDay, 0, len(req.Days))

	for _, d := range req.Days {
		day := &entity.ScheduleDay{
			ScheduleID: schedule.ID,
			Weekday:    d.Weekday,
			DayIndex:   d.DayIndex,
			IsDayOff:   d.IsDayOff,
		}
		if d.StartTime != nil && d.EndTime != nil {
			start, end, err := parseTimeDuration(*d.StartTime, *d.EndTime)
			if err != nil {
				return errors.NewHTTPError(http.StatusBadRequest, "invalid time range", err)
			}
			day.StartTime = &start
			day.EndTime = &end
		}

		days = append(days, day)
	}

	schedule, err = h.scheduleUseCase.CreateSchedule(ctx, schedule, days)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
	}

	log.Info("schedule created", "schedule_id", schedule.ID)
	return c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid ID", err)
	}

	schedule, err := h.scheduleUseCase.GetScheduleByID(ctx, id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to get schedule", err)
	}

	log.Info("schedule retrieved", "schedule_id", id)
	return c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) ListSchedules(c echo.Context) error {
	ctx := c.Request().Context()

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 100
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	if limit > 1000 {
		limit = 1000
	}

	schedules, err := h.scheduleUseCase.ListSchedules(ctx, offset, limit)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to list schedules", err)
	}
	return c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) UpdateSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid ID", err)
	}

	var req dto.UpdateScheduleRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid body", err)
	}

	scheduleType, err := req.Type.ToEntity()
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid type value", err)
	}

	schedule := &entity.Schedule{
		ID:        id,
		MasterID:  req.MasterID,
		Name:      req.Name,
		Type:      scheduleType,
		StartDate: time.Now(),
		EndDate:   req.EndDate,
	}
	if req.StartDate != nil {
		schedule.StartDate = *req.StartDate
	}

	if err := h.scheduleUseCase.UpdateSchedule(ctx, schedule); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to update schedule", err)
	}

	log.Info("schedule updated", "schedule_id", id)
	return c.NoContent(http.StatusOK)
}

func (h *ScheduleHandler) DeleteSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.scheduleUseCase.DeleteSchedule(ctx, id); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to delete schedule", err)
	}

	log.Info("schedule deleted", "schedule_id", id)
	return c.NoContent(http.StatusNoContent)
}

func (h *ScheduleHandler) AddDay(c echo.Context) error {
	ctx := c.Request().Context()

	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid schedule id", err)
	}

	var req dto.CreateScheduleDayRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid body", err)
	}

	day := &entity.ScheduleDay{
		ScheduleID: scheduleID,
		Weekday:    req.Weekday,
		DayIndex:   req.DayIndex,
		IsDayOff:   req.IsDayOff,
	}
	if req.StartTime != nil && req.EndTime != nil {
		start, end, err := parseTimeDuration(*req.StartTime, *req.EndTime)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid time range", err)
		}
		day.StartTime = &start
		day.EndTime = &end
	}

	if err := h.scheduleUseCase.AddDay(ctx, day); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to add day", err)
	}
	return c.JSON(http.StatusCreated, day)
}

// GET /api/v1/schedules/:id/days
func (h *ScheduleHandler) ListDays(c echo.Context) error {
	ctx := c.Request().Context()

	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid schedule id", err)
	}

	days, err := h.scheduleUseCase.GetDaysBySchedule(ctx, scheduleID)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list days", err)
	}
	return c.JSON(http.StatusOK, days)
}

// PUT /api/v1/schedules/days/:id
func (h *ScheduleHandler) UpdateDay(c echo.Context) error {
	ctx := c.Request().Context()

	dayID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid day id", err)
	}

	var req dto.UpdateScheduleDayRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid body", err)
	}

	day := &entity.ScheduleDay{
		ID:       dayID,
		Weekday:  req.Weekday,
		DayIndex: req.DayIndex,
		IsDayOff: req.IsDayOff,
	}
	if req.StartTime != nil && req.EndTime != nil {
		start, end, err := parseTimeDuration(*req.StartTime, *req.EndTime)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid time range", err)
		}
		day.StartTime = &start
		day.EndTime = &end
	}

	if err := h.scheduleUseCase.UpdateDay(ctx, day); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update day", err)
	}
	return c.NoContent(http.StatusOK)
}

// DELETE /api/v1/schedules/days/:day_id
func (h *ScheduleHandler) DeleteDay(c echo.Context) error {
	ctx := c.Request().Context()

	dayID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid day id", err)
	}

	if err := h.scheduleUseCase.DeleteDay(ctx, dayID); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete day", err)
	}
	return c.NoContent(http.StatusNoContent)
}

// POST /api/v1/schedules/:id/slots
func (h *ScheduleHandler) AddSlot(c echo.Context) error {
	ctx := c.Request().Context()

	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid schedule id", err)
	}

	var req dto.CreateScheduleSlotRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid body", err)
	}

	slot := &entity.ScheduleSlot{
		ScheduleID: scheduleID,
		Date:       timeutil.NormalizeDate(req.Date),
		IsDayOff:   req.IsDayOff,
	}

	if req.StartTime != nil && req.EndTime != nil {
		start, end, err := parseTimeDuration(*req.StartTime, *req.EndTime)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid time range", err)
		}
		slot.StartTime = &start
		slot.EndTime = &end
	}

	if err := h.scheduleUseCase.AddSlot(ctx, slot); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to add slot", err)
	}
	return c.JSON(http.StatusCreated, slot)
}

// GET /api/v1/schedules/:id/slots
func (h *ScheduleHandler) ListSlots(c echo.Context) error {
	ctx := c.Request().Context()

	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid schedule id", err)
	}

	slots, err := h.scheduleUseCase.GetSlotsBySchedule(ctx, scheduleID)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list slots", err)
	}
	return c.JSON(http.StatusOK, slots)
}

// PUT /api/v1/schedules/slots/:id
func (h *ScheduleHandler) UpdateSlot(c echo.Context) error {
	ctx := c.Request().Context()

	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid slot id", err)
	}

	var req dto.UpdateScheduleSlotRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid body", err)
	}

	slot := &entity.ScheduleSlot{
		ID:       slotID,
		Date:     timeutil.NormalizeDate(req.Date),
		IsDayOff: req.IsDayOff,
	}
	if req.StartTime != nil && req.EndTime != nil {
		start, end, err := parseTimeDuration(*req.StartTime, *req.EndTime)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid time range", err)
		}
		slot.StartTime = &start
		slot.EndTime = &end
	}

	if err := h.scheduleUseCase.UpdateSlot(ctx, slot); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update slot", err)
	}
	return c.NoContent(http.StatusOK)
}

// DELETE /api/v1/schedules/slots/:id
func (h *ScheduleHandler) DeleteSlot(c echo.Context) error {
	ctx := c.Request().Context()

	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid slot id", err)
	}

	if err := h.scheduleUseCase.DeleteSlot(ctx, slotID); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete slot", err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GET /api/v1/schedules/master/:master_id/date/:date
func (h *ScheduleHandler) GetScheduleForDate(c echo.Context) error {
	ctx := c.Request().Context()

	masterID, err := uuid.Parse(c.Param("master_id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid master id", err)
	}

	dateStr := c.Param("date")
	if dateStr == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "date is required", nil)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid date format (use YYYY-MM-DD)", err)
	}
	date = timeutil.NormalizeDate(date)

	resp, err := h.scheduleUseCase.GetScheduleForDate(ctx, masterID, date)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to get schedule for date", err)
	}

	// Если ничего нет — возвращаем пустой массив, чтобы клиенту было проще
	if resp == nil {
		resp = []dto.ScheduleForDateResponse{}
	}

	return c.JSON(http.StatusOK, resp)
}

// GET /api/v1/schedules/master/:master_id/range?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *ScheduleHandler) GetScheduleForRange(c echo.Context) error {
	ctx := c.Request().Context()

	masterID, err := uuid.Parse(c.Param("master_id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid master id", err)
	}

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	if fromStr == "" || toStr == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "from and to params are required", nil)
	}

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid from date format (use YYYY-MM-DD)", err)
	}
	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid to date format (use YYYY-MM-DD)", err)
	}

	if toDate.Before(fromDate) {
		return errors.NewHTTPError(http.StatusBadRequest, "toDate must be greater than fromDate", nil)
	}

	fromDate = timeutil.NormalizeDate(fromDate)
	toDate = timeutil.NormalizeDate(toDate)

	resp, err := h.scheduleUseCase.GetScheduleForRange(ctx, masterID, fromDate, toDate)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to get schedule range", err)
	}

	return c.JSON(http.StatusOK, resp)
}

// ---------- helpers ----------

func parseTimeOfDay(s string) (time.Time, error) {
	// ожидаем формат "15:04"
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func parseTimeDuration(from, to string) (time.Time, time.Time, error) {
	// парсим "HH:MM"
	start, err := parseTimeOfDay(from)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parseTimeOfDay(to)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if !start.Before(end) {
		return time.Time{}, time.Time{}, errors.ErrEndTimeBeforeStartTime
	}

	return start, end, nil
}
