package application

import (
	"context"
	"fmt"

	"github.com/go-brick-template/go-brick-template/internal/application/runtime"
	"github.com/go-brick-template/go-brick-template/internal/config"
	adminhttp "github.com/go-brick-template/go-brick-template/module/admin_http"
	httpmod "github.com/go-brick-template/go-brick-template/module/http"
	"github.com/go-brick-template/go-brick-template/module/postgres"
)

func Run(ctx context.Context) error {
	cfg := config.LoadFromEnv()
	if cfg.Postgres == nil {
		return fmt.Errorf("DATABASE_URL is required")
	}

	pgMod, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return err
	}

	reg := wireBricks(pgMod.DB())
	httpMod := httpmod.New(cfg.HTTPServer)
	adminMod := adminhttp.New(cfg.AdminServer, pgMod.DB())

	app := runtime.New("go-brick-template").
		WithModule(pgMod).
		WithModule(httpMod).
		WithModule(adminMod)

	if err := app.Build(ctx); err != nil {
		return err
	}

	httpMod.MountRoutes(cfg, reg)

	c := runtime.NewCloser()
	app.Run(ctx, c)
	return c.Close()
}
