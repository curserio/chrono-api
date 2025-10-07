package entity

import (
	"time"

	"github.com/google/uuid"
)

type Master struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            *string   `json:"email,omitempty"`
	Phone            *string   `json:"phone,omitempty"`
	TelegramID       *int64    `json:"telegram_id,omitempty"`
	TelegramUsername *string   `json:"telegram_username,omitempty"`
	Description      *string   `json:"description,omitempty"`
	City             *string   `json:"city,omitempty"`
	Timezone         string    `json:"timezone"`
	Language         string    `json:"language"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
