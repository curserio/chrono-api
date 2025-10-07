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
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MasterHandler struct {
	masterUseCase *usecase.MasterUseCase
}

func NewMasterHandler(s *server.Server, uc *usecase.MasterUseCase) {
	handler := &MasterHandler{
		masterUseCase: uc,
	}

	// Routes
	group := s.NewGroup("/api/v1/masters")
	group.POST("", handler.CreateMaster)
	group.GET("", handler.ListMasters)
	group.GET("/:id", handler.GetMaster)
	group.PUT("/:id", handler.UpdateMaster)
	group.DELETE("/:id", handler.DeleteMaster)
}

func (h *MasterHandler) CreateMaster(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	req := new(dto.CreateMasterRequest)
	if err := c.Bind(req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
	}

	if err := c.Validate(req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Validation failed", err)
	}

	master := &entity.Master{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		TelegramID:       req.TelegramID,
		TelegramUsername: req.TelegramUsername,
		City:             req.City,
		Timezone:         req.Timezone,
		Language:         req.Language,
		Description:      req.Description,
	}

	master, err := h.masterUseCase.CreateMaster(ctx, master)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
	}

	log.Info("master created", "master_id", master.ID)
	return c.JSON(http.StatusCreated, master)
}

func (h *MasterHandler) GetMaster(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid ID", err)
	}

	master, err := h.masterUseCase.GetMasterByID(ctx, id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
	}

	log.Info("master retrieved", "master_id", id)
	return c.JSON(http.StatusOK, master)
}

// GET /masters
func (h *MasterHandler) ListMasters(c echo.Context) error {
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

	masters, err := h.masterUseCase.ListMasters(ctx, offset, limit)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list masters", err)
	}
	return c.JSON(http.StatusOK, masters)
}

// PUT /masters/:id
func (h *MasterHandler) UpdateMaster(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	var req dto.UpdateMasterRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	master := &entity.Master{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		TelegramID:       req.TelegramID,
		TelegramUsername: req.TelegramUsername,
		Description:      req.Description,
		City:             req.City,
		Timezone:         req.Timezone,
		Language:         req.Language,
		UpdatedAt:        time.Now(),
	}
	master.ID = id

	if err := h.masterUseCase.UpdateMaster(ctx, master); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update master", err)
	}

	log.Info("master updated", "master_id", id)
	return c.NoContent(http.StatusOK)
}

// DELETE /masters/:id
func (h *MasterHandler) DeleteMaster(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	if err := h.masterUseCase.DeleteMaster(ctx, id); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete master", err)
	}

	log.Info("master deleted", "master_id", id)
	return c.NoContent(http.StatusNoContent)
}
