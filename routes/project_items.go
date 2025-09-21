package routes


import (
	"github.com/clem-kay/mini-trello/middleware"
	"github.com/clem-kay/mini-trello/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterProjectItemsRoutes(app *fiber.App) {
	api := app.Group("api/v1/items")

	api.Post("/", middleware.AuthMiddleware(), middleware.ValidateBody(), services.CreateProjectItem)
	api.Get("/", services.GetProjectItems)
	api.Get("/:id", services.GetProjectItemByID)
	api.Put("/:id", services.UpdateProjectItem)
	api.Get("/board/:id", services.GetProjectItemsByBoardID)
	api.Delete("/:id", services.DeleteProjectItem)

}
