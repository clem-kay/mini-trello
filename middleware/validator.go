// middleware/validator.go
package middleware

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/clem-kay/mini-trello/utils"
)

var validate = validator.New()

func ValidateBody() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Locals("body") // or parse directly
		if err := validate.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}
		return c.Next()
	}
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		//Token validation logic here. take the logicfrom the auth header
		bearerToken:= strings.Split(authHeader, "")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			return c. Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected 'bearer <token>'",
			})
		}

		tokenString := bearerToken[1]

		userID, err := utils.GetUserIDFromToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token" + err.Error(),
				// Optional: for debugging — remove in prod
				// "details": err.Error(),
			})
		}

		// ✅ Attach to context
		c.Locals("user_id", userID)

		// Optional: Also extract and attach email if needed
		// (You can create a similar helper: GetEmailFromToken)

		return c.Next()
	}
}
