package tool

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
)

type httpError struct {
	Error string `json:"error"`
}

// Maps core errors to HTTP ones and sends them to the client
// If an error can't be mapped, it sends 500 status code by default and logs the error
func MapHttpError(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var status int
	var text string

	switch {
	case errors.Is(err, db.ErrEntityNotFound):
		status = 404
	case errors.Is(err, fiber.ErrUnprocessableEntity):
		status = 400
	default:
		status = 500
		log.Printf(
			"ERROR: %s - %s %s %s, status: %d, %s\n",
			c.IP(),
			c.Method(),
			c.OriginalURL(),
			c.Protocol(),
			status,
			err.Error(),
		)
	}

	if text == "" {
		return c.Status(status).End()
	}
	return c.Status(status).JSON(httpError{
		Error: err.Error(),
	})
}
