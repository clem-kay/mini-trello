package routes

import (
	"github.com/clem-kay/mini-trello/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthorRoutes(app *fiber.App) {
	api := app.Group("api/v1/auth")

	api.Post("/login", services.Login)
	api.Post("/register", services.Register)

}
