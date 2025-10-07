package entity

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	TelegramID       *int64    `json:"telegram_id,omitempty"`
	TelegramUsername *string   `json:"telegram_username,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
