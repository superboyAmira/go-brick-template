package http

import (
	"context"
	"log/slog"

	"github.com/go-brick-template/go-brick-template/internal/api/handlers"
	"github.com/go-brick-template/go-brick-template/internal/registry"
	"github.com/go-brick-template/go-brick-template/internal/application/runtime"
	"github.com/go-brick-template/go-brick-template/internal/config"
	"github.com/go-brick-template/go-brick-template/internal/config/options"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Module struct {
	cfg *options.HTTPOptions
	app *fiber.App
}

func New(cfg *options.HTTPOptions) *Module {
	return &Module{cfg: cfg}
}

func (m *Module) Init(_ context.Context, _ runtime.Info) error {
	if m.cfg == nil {
		return nil
	}
	m.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             32 * 1024 * 1024,
	})
	m.app.Use(recover.New())
	return nil
}

func (m *Module) MountRoutes(cfg *config.Config, reg *registry.Registry) {
	if m.app == nil || reg == nil {
		return
	}
	handlers.Register(m.app, cfg, reg)
}

func (m *Module) Run(_ context.Context, c *runtime.Closer) {
	if m.app == nil || m.cfg == nil {
		return
	}
	addr := m.cfg.Addr
	go func() {
		slog.Info("http listening", "addr", addr)
		if err := m.app.Listen(addr); err != nil {
			slog.Error("http server stopped", "error", err)
		}
	}()
	c.Add("http", func() error {
		if m.app == nil {
			return nil
		}
		return m.app.Shutdown()
	})
}
