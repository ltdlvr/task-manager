package rest

import "github.com/gofiber/fiber/v3"

const userIDLocalKey = "userID"

// gets userID from Locals in fiber.Ctx
func getUserID(c fiber.Ctx) (uint64, error) {
	userID, ok := c.Locals(userIDLocalKey).(uint64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return userID, nil
}
