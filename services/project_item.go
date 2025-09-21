package services

import (
	"time"

	"github.com/clem-kay/mini-trello/config"
	"github.com/clem-kay/mini-trello/models"
	"github.com/clem-kay/mini-trello/utils"
	"github.com/gofiber/fiber/v2"
)

type ItemRequestPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // ISO 8601 format
	Status      string `json:"status" validate:"oneof=todo in_progress done"`
	Priority    string `json:"priority" validate:"oneof=low medium high"`
	BoardID     uint   `json:"board_id" validate:"required"`
}

// ✅ Create Project Item
func CreateProjectItem(c *fiber.Ctx) error {
	currentUserID := utils.GetIDFromContext(c)
	if currentUserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, user needs to be logged in",
		})
	}

	var body ItemRequestPayload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Name == "" || body.BoardID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and BoardID are required",
		})
	}

	// Check if board exists
	var board models.Board
	result := config.DB.First(&board, body.BoardID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Board not found",
		})
	}

	// Parse due date
	var dueDate *time.Time
	if body.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, body.DueDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid due_date format, must be RFC3339 (e.g. 2025-09-20T15:04:05Z)",
			})
		}
		dueDate = &parsed
	}

	item := models.ProjectItem{
		Name:        body.Name,
		BoardID:     body.BoardID,
		Description: body.Description,
		Status:      models.ItemStatus(body.Status),
		Priority:    models.ItemPriority(body.Priority),
		DueDate:     dueDate,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create item: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Item created successfully",
		"item":    item,
	})
}

// ✅ Get all project items
func GetProjectItems(c *fiber.Ctx) error {
	var items []models.ProjectItem
	if err := config.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch items: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

// ✅ Get single item by ID
func GetProjectItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.ProjectItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item fetched successfully",
		"item":    item,
	})
}

// ✅ Update item
func UpdateProjectItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.ProjectItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	var body ItemRequestPayload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	item.Name = body.Name
	item.Description = body.Description
	item.Status = models.ItemStatus(body.Status)
	item.Priority = models.ItemPriority(body.Priority)

	if body.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, body.DueDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid due_date format",
			})
		}
		item.DueDate = &parsed
	}

	if err := config.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update item: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item updated successfully",
		"item":    item,
	})
}

func DeleteProjectItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.ProjectItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	if err := config.DB.Delete(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete item: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Item deleted successfully",
	})
}

// ✅ Get items by board ID
func GetProjectItemsByBoardID(c *fiber.Ctx) error {
	boardID := c.Params("id")
	var items []models.ProjectItem

	if err := config.DB.Where("board_id = ?", boardID).Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch items: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Items fetched successfully",
		"items":   items,
	})

}
