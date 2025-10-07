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

type ClientHandler struct {
	clientUseCase *usecase.ClientUseCase
}

func NewClientHandler(s *server.Server, uc *usecase.ClientUseCase) {
	handler := &ClientHandler{clientUseCase: uc}

	group := s.NewGroup("/api/v1/clients")
	group.POST("", handler.CreateClient)
	group.GET("/:id", handler.GetClient)
	group.GET("", handler.ListClients)
	group.PUT("/:id", handler.UpdateClient)
	group.DELETE("/:id", handler.DeleteClient)
}

func (h *ClientHandler) CreateClient(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	var req dto.CreateClientRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	client, err := h.clientUseCase.CreateClient(ctx, &entity.Client{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		TelegramID:       req.TelegramID,
		TelegramUsername: req.TelegramUsername,
		City:             req.City,
		Timezone:         req.Timezone,
		Language:         req.Language,
	})
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to create client", err)
	}

	log.Info("client created", "client_id", client.ID)
	return c.JSON(http.StatusCreated, client)
}

func (h *ClientHandler) GetClient(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	client, err := h.clientUseCase.GetClientByID(ctx, id)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to get client", err)
	}

	return c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) ListClients(c echo.Context) error {
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

	clients, err := h.clientUseCase.ListClients(ctx, offset, limit)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to list clients", err)
	}
	return c.JSON(http.StatusOK, clients)
}

func (h *ClientHandler) UpdateClient(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	var req dto.UpdateClientRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid request", err)
	}

	client := &entity.Client{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		TelegramID:       req.TelegramID,
		TelegramUsername: req.TelegramUsername,
		City:             req.City,
		Timezone:         req.Timezone,
		Language:         req.Language,
		UpdatedAt:        time.Now(),
	}
	client.ID = id
	if err := h.clientUseCase.UpdateClient(ctx, client); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to update client", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *ClientHandler) DeleteClient(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "invalid id", err)
	}

	if err := h.clientUseCase.DeleteClient(ctx, id); err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "failed to delete client", err)
	}
	return c.NoContent(http.StatusNoContent)
}
