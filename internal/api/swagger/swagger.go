// Package swagger serves Swagger UI from the canonical OpenAPI contract
// (docs/contracts/openapi.yaml). Regenerate embed: make swagger-gen
package swagger

import (
	_ "embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi.yaml
var openAPIYAML []byte

//go:embed admin-openapi.yaml
var adminOpenAPIYAML []byte

//go:embed index.html
var indexHTML []byte

//go:embed admin.html
var adminIndexHTML []byte

// Register mounts OpenAPI specs and Swagger UI on the main API server (public, no JWT).
func Register(app *fiber.App) {
	app.Get("/openapi.yaml", serveYAML(openAPIYAML))
	app.Get("/swagger", serveHTML(indexHTML))
	app.Get("/swagger/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger", fiber.StatusMovedPermanently)
	})
}

// RegisterAdmin mounts admin OpenAPI + Swagger UI on the admin HTTP server.
func RegisterAdmin(app *fiber.App) {
	app.Get("/openapi.yaml", serveYAML(adminOpenAPIYAML))
	app.Get("/swagger", serveHTML(adminIndexHTML))
	app.Get("/swagger/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger", fiber.StatusMovedPermanently)
	})
}

func serveYAML(spec []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
		c.Set(fiber.HeaderCacheControl, "no-cache")
		return c.Status(http.StatusOK).Send(spec)
	}
}

func serveHTML(page []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		c.Set(fiber.HeaderCacheControl, "no-cache")
		return c.Status(http.StatusOK).Send(page)
	}
}
