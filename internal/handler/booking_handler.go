package handler

import (
	"net/http"
	"time"

	"github.com/curserio/chrono-api/internal/domain/entity"
	"github.com/curserio/chrono-api/internal/dto"
	"github.com/curserio/chrono-api/internal/errors"
	"github.com/curserio/chrono-api/internal/infrastructure/http/server"
	"github.com/curserio/chrono-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingUseCase *usecase.BookingUseCase
}

func NewBookingHandler(s *server.Server, uc *usecase.BookingUseCase) {
	handler := &BookingHandler{bookingUseCase: uc}

	group := s.NewGroup("/api/v1/bookings")
	group.POST("", handler.CreateBooking)
	group.GET("/:id", handler.GetBooking)
	group.GET("/master/:master_id", handler.GetByMaster)
	group.GET("/client/:client_id", handler.GetByClient)
	group.PUT("/:id/status", handler.UpdateStatus)
	group.DELETE("/:id", handler.DeleteBooking)
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	start, end, err := parseTimeDuration(req.StartTime, req.EndTime)
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid time format", err)
	}

	// TODO сохраняю в UTC, хотя у мастера может быть другая таймзона
	date := req.Date.Truncate(24 * time.Hour).UTC()
	start = time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.UTC)
	end = time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.UTC)

	booking, err := h.bookingUseCase.CreateBooking(ctx, &entity.Booking{
		MasterID:  req.MasterID,
		ClientID:  req.ClientID,
		ServiceID: req.ServiceID,
		StartTime: start,
		EndTime:   end,
		Status:    entity.BookingStatusPending,
	})
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to create booking", err)
	}
	return c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) GetBooking(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}
	b, err := h.bookingUseCase.GetBookingByID(ctx, id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to get booking", err)
	}
	return c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) GetByMaster(c echo.Context) error {
	ctx := c.Request().Context()
	masterID, err := uuid.Parse(c.Param("master_id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid master id", err)
	}

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	var from, to time.Time

	if fromStr != "" {
		from, err = time.Parse(time.DateOnly, fromStr)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid 'from' date format (use YYYY-MM-DD)", err)
		}
	} else {
		from = time.Now().Truncate(24 * time.Hour)
	}

	if toStr != "" {
		to, err = time.Parse(time.DateOnly, toStr)
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, "invalid 'to' date format (use YYYY-MM-DD)", err)
		}
	} else {
		to = from.AddDate(0, 0, 7)
	}

	if to.Before(from) {
		return errors.NewHTTPError(http.StatusBadRequest, "'to' date must be after 'from'", nil)
	}

	bookings, err := h.bookingUseCase.GetBookingsByMaster(ctx, masterID, from, to)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list bookings", err)
	}
	return c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetByClient(c echo.Context) error {
	ctx := c.Request().Context()
	clientID, err := uuid.Parse(c.Param("client_id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid client id", err)
	}
	bookings, err := h.bookingUseCase.GetBookingsByClient(ctx, clientID)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list bookings", err)
	}
	return c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) UpdateStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}
	var req dto.UpdateBookingStatusRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	status, err := req.Status.ToEntity()
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid status value", err)
	}

	if err := h.bookingUseCase.UpdateBookingStatus(ctx, id, status); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update status", err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BookingHandler) DeleteBooking(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}
	if err := h.bookingUseCase.DeleteBooking(ctx, id); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete booking", err)
	}
	return c.NoContent(http.StatusNoContent)
}
