package application

import (
	itemrepo "github.com/go-brick-template/go-brick-template/internal/brick/item/repository"
	itemsvc "github.com/go-brick-template/go-brick-template/internal/brick/item/service"
	"github.com/go-brick-template/go-brick-template/internal/registry"
	"github.com/go-brick-template/go-brick-template/module/postgres"
)

func wireBricks(pg *postgres.DB) *registry.Registry {
	itemRepo := itemrepo.NewPostgres(pg)
	itemSvc := itemsvc.New(itemRepo)

	return &registry.Registry{
		Item: itemSvc,
	}
}
