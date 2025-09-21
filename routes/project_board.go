package routes

import (
	"github.com/clem-kay/mini-trello/middleware"
	"github.com/clem-kay/mini-trello/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterBoardoutes(app *fiber.App) {
	api := app.Group("api/v1/boards")

	api.Post("/", middleware.AuthMiddleware(), middleware.ValidateBody(), services.CreateBoard)
	api.Get("/", services.GetBoards)
	api.Get("/:id", services.GetBoardByID)
	api.Put("/:id", services.UpdateBoard)
	api.Get("/user/:id", services.GetBoardByUserID)
	api.Delete("/:id", services.DeleteBoard)

}
