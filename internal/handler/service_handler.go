package handler

import (
	"net/http"
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

type ServiceHandler struct {
	serviceUseCase *usecase.ServiceUseCase
}

func NewServiceHandler(s *server.Server, uc *usecase.ServiceUseCase) {
	handler := &ServiceHandler{serviceUseCase: uc}

	group := s.NewGroup("/api/v1/services")
	group.POST("", handler.CreateService)
	group.GET("/:id", handler.GetService)
	group.GET("/master/:master_id", handler.ListByMaster)
	group.PUT("/:id", handler.UpdateService)
	group.DELETE("/:id", handler.DeleteService)
}

func (h *ServiceHandler) CreateService(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	var req dto.CreateServiceRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	service, err := h.serviceUseCase.CreateService(ctx, &entity.Service{
		MasterID:    req.MasterID,
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
	})
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to create service", err)
	}

	log.Info("service created", "service_id", service.ID)
	return c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) GetService(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}
	svc, err := h.serviceUseCase.GetServiceByID(ctx, id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to get service", err)
	}
	return c.JSON(http.StatusOK, svc)
}

func (h *ServiceHandler) ListByMaster(c echo.Context) error {
	ctx := c.Request().Context()
	masterID, err := uuid.Parse(c.Param("master_id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid master id", err)
	}

	svcs, err := h.serviceUseCase.ListServicesByMaster(ctx, masterID)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list services", err)
	}
	return c.JSON(http.StatusOK, svcs)
}

func (h *ServiceHandler) UpdateService(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	var req dto.UpdateServiceRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	service := &entity.Service{
		MasterID:    req.MasterID,
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
		UpdatedAt:   time.Now(),
	}
	service.ID = id
	if err := h.serviceUseCase.UpdateService(ctx, service); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update service", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	if err := h.serviceUseCase.DeleteService(ctx, id); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete service", err)
	}
	return c.NoContent(http.StatusNoContent)
}
