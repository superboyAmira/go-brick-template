package adminhttp

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-brick-template/go-brick-template/internal/api/swagger"
	"github.com/go-brick-template/go-brick-template/internal/application/runtime"
	"github.com/go-brick-template/go-brick-template/internal/config/options"
	"github.com/go-brick-template/go-brick-template/module/postgres"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	cfg *options.HTTPOptions
	pg  *postgres.DB
	app *fiber.App
}

func New(cfg *options.HTTPOptions, pg *postgres.DB) *Module {
	return &Module{cfg: cfg, pg: pg}
}

func (m *Module) Init(ctx context.Context, _ runtime.Info) error {
	if m.cfg == nil {
		return nil
	}
	m.app = fiber.New(fiber.Config{DisableStartupMessage: true})
	swagger.RegisterAdmin(m.app)
	m.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	adminToken := m.cfg.Token
	protected := m.app.Group("", adminTokenMiddleware(adminToken))
	protected.Get("/ready", func(c *fiber.Ctx) error {
		if err := postgres.Ping(ctx, m.pg); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"error":  err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})
	protected.Get("/metrics", func(c *fiber.Ctx) error {
		return c.SendString("# go-brick-template metrics placeholder\n")
	})
	return nil
}

func adminTokenMiddleware(token string) fiber.Handler {
	if token == "" {
		return func(c *fiber.Ctx) error { return c.Next() }
	}
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "Bearer "+token || c.Get("X-Admin-Token") == token {
			return c.Next()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
}

func (m *Module) Run(ctx context.Context, c *runtime.Closer) {
	if m.app == nil || m.cfg == nil {
		return
	}
	addr := m.cfg.Addr
	go func() {
		slog.Info("admin http listening", "addr", addr)
		if err := m.app.Listen(addr); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("admin http stopped", "error", err)
		}
	}()
	c.Add("admin_http", func() error {
		if m.app == nil {
			return nil
		}
		return m.app.Shutdown()
	})
}
