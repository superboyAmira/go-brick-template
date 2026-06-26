package model

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
