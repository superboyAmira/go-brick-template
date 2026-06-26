package service

import (
	"context"
	"errors"

	"github.com/go-brick-template/go-brick-template/internal/brick/item"
	"github.com/go-brick-template/go-brick-template/internal/brick/item/model"
	"github.com/go-brick-template/go-brick-template/internal/brick/item/repository"
	apperr "github.com/go-brick-template/go-brick-template/internal/shared/apperr"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *repository.PostgresRepository
}

func New(repo *repository.PostgresRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]item.Item, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return mapItems(items), nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*item.Item, error) {
	it, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrItemNotFound
		}
		return nil, err
	}
	out := mapItem(*it)
	return &out, nil
}

func (s *Service) Create(ctx context.Context, req item.CreateRequest) (*item.Item, error) {
	if req.Title == "" {
		return nil, apperr.Validation("title is required", nil)
	}
	it, err := s.repo.Create(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	out := mapItem(*it)
	return &out, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.ErrItemNotFound
		}
		return err
	}
	return nil
}

func mapItem(it model.Item) item.Item {
	return item.Item{
		ID:        it.ID,
		Title:     it.Title,
		CreatedAt: it.CreatedAt,
		UpdatedAt: it.UpdatedAt,
	}
}

func mapItems(items []model.Item) []item.Item {
	out := make([]item.Item, 0, len(items))
	for _, it := range items {
		out = append(out, mapItem(it))
	}
	return out
}

var _ item.Service = (*Service)(nil)
