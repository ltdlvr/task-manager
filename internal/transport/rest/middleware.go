package rest

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	adapterTool "github.com/ltdlvr/task-manager/internal/core/adapter/tool"
)

func AuthMiddleware(tokenTool adapterTool.Token) fiber.Handler {
	return func(c fiber.Ctx) error {

		header := c.Get("Authorization")
		if header == "" {
			return fiber.ErrUnauthorized
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenString == "" {
			return fiber.ErrUnauthorized
		}

		userID, err := tokenTool.Parse(tokenString)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals(userIDLocalKey, userID)

		return c.Next()
	}
}
