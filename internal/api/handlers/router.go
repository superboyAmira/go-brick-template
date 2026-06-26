package handlers

import (
	"github.com/go-brick-template/go-brick-template/internal/api/middleware"
	"github.com/go-brick-template/go-brick-template/internal/api/response"
	"github.com/go-brick-template/go-brick-template/internal/api/swagger"
	"github.com/go-brick-template/go-brick-template/internal/brick/item"
	"github.com/go-brick-template/go-brick-template/internal/config"
	"github.com/go-brick-template/go-brick-template/internal/registry"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Register(app *fiber.App, _ *config.Config, reg *registry.Registry) {
	swagger.Register(app)

	api := app.Group("/api/v1")
	api.Use(middleware.RequestID())
	registerItemRoutes(api, reg)
}

func registerItemRoutes(api fiber.Router, reg *registry.Registry) {
	api.Get("/items", listItems(reg))
	api.Post("/items", createItem(reg))
	api.Get("/items/:id", getItem(reg))
	api.Delete("/items/:id", deleteItem(reg))
}

func listItems(reg *registry.Registry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, err := reg.Item.List(c.Context())
		if err != nil {
			return response.WriteError(c, err)
		}
		return response.WriteJSON(c, fiber.StatusOK, fiber.Map{"items": items})
	}
}

func createItem(reg *registry.Registry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req item.CreateRequest
		if err := c.BodyParser(&req); err != nil {
			return response.WriteError(c, err)
		}
		it, err := reg.Item.Create(c.Context(), req)
		if err != nil {
			return response.WriteError(c, err)
		}
		return response.WriteJSON(c, fiber.StatusCreated, it)
	}
}

func getItem(reg *registry.Registry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return response.WriteError(c, err)
		}
		it, err := reg.Item.Get(c.Context(), id)
		if err != nil {
			return response.WriteError(c, err)
		}
		return response.WriteJSON(c, fiber.StatusOK, it)
	}
}

func deleteItem(reg *registry.Registry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return response.WriteError(c, err)
		}
		if err := reg.Item.Delete(c.Context(), id); err != nil {
			return response.WriteError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
