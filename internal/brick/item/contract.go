package item

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRequest struct {
	Title string `json:"title"`
}

type Service interface {
	List(ctx context.Context) ([]Item, error)
	Get(ctx context.Context, id uuid.UUID) (*Item, error)
	Create(ctx context.Context, req CreateRequest) (*Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
