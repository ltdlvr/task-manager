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
	var message string

	switch {
	case errors.Is(err, fiber.ErrBadRequest):
		status = fiber.StatusBadRequest
		message = "bad request"
	case errors.Is(err, fiber.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		message = "unauthorized"
	case errors.Is(err, fiber.ErrForbidden):
		status = fiber.StatusForbidden
		message = "forbidden"
	case errors.Is(err, db.ErrForbidden):
		status = fiber.StatusForbidden
		message = "forbidden"
	case errors.Is(err, db.ErrEntityNotFound):
		status = fiber.StatusNotFound
		message = "not found"
	case errors.Is(err, fiber.ErrUnprocessableEntity):
		status = fiber.StatusUnprocessableEntity
		message = "unprocessable entity"
	case errors.Is(err, db.ErrOutOfRange):
		status = fiber.StatusUnprocessableEntity
		message = "unprocessable entity"
	default:
		status = fiber.StatusInternalServerError
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

	if message == "" {
		return c.SendStatus(status)
	}
	return c.Status(status).JSON(httpError{
		Error: message,
	})
}
