package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID `json:"id"`
	MasterID    uuid.UUID `json:"master_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"` // in minutes
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
