package response

import (
	apperr "github.com/go-brick-template/go-brick-template/internal/shared/apperr"

	"github.com/gofiber/fiber/v2"
)

type ErrorBody struct {
	Error ErrorEnvelope `json:"error"`
}

type ErrorEnvelope struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
	RequestID string         `json:"request_id"`
}

func WriteError(c *fiber.Ctx, err error) error {
	ae := apperr.From(err)
	rid, _ := c.Locals("request_id").(string)
	return c.Status(ae.HTTPStatus).JSON(ErrorBody{
		Error: ErrorEnvelope{
			Code:      ae.Code,
			Message:   ae.Message,
			Details:   ae.Details,
			RequestID: rid,
		},
	})
}

func WriteJSON(c *fiber.Ctx, status int, body any) error {
	return c.Status(status).JSON(body)
}

type CursorPage struct {
	Items      any     `json:"items"`
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}
