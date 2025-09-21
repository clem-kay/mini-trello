// services/board.go
package services

import (
	"strconv"

	"github.com/clem-kay/mini-trello/config"
	"github.com/clem-kay/mini-trello/models"
	"github.com/clem-kay/mini-trello/utils"
	"github.com/gofiber/fiber/v2"
)

type boardRequest struct {
	Name        string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
}

// ✅ CREATE
func CreateBoard(c *fiber.Ctx) error {
	currentUserID := utils.GetIDFromContext(c)
	if currentUserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, user needs to be logged in",
		})
	}

	var body boardRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	board := models.Board{
		Name:        body.Name,
		Description: body.Description,
		UserID:      currentUserID,
	}

	if err := config.DB.Create(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create board: " + err.Error(), // ← Fixed: added space
		})
	}

	return c.Status(fiber.StatusCreated).JSON(board)
}

// ✅ GET ALL (public for now, or you can protect it)
func GetBoards(c *fiber.Ctx) error {
	var boardList []models.Board

	// Optional: add pagination, limit, etc.
	result := config.DB.Find(&boardList)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch boards",
		})
	}

	return c.JSON(boardList)
}

// ✅ GET BY ID
func GetBoardByID(c *fiber.Ctx) error {
	idStr := utils.GetIDParam(c)
	if idStr == "" {
		return nil // error already sent
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var board models.Board
	if err := config.DB.First(&board, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Board not found",
		})
	}

	return c.JSON(board)
}

// ✅ UPDATE (only owner can update)
func UpdateBoard(c *fiber.Ctx) error {
	currentUserID := utils.GetIDFromContext(c)
	if currentUserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, user needs to be logged in",
		})
	}

	idStr := utils.GetIDParam(c)
	if idStr == "" {
		return nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var board models.Board
	if err := config.DB.First(&board, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Board not found",
		})
	}

	// ✅ Ownership check
	if board.UserID != currentUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not own this board",
		})
	}

	var body boardRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	// Update fields
	board.Name = body.Name
	board.Description = body.Description

	if err := config.DB.Save(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update board: " + err.Error(),
		})
	}

	return c.JSON(board)
}

// ✅ DELETE (only owner can delete)
func DeleteBoard(c *fiber.Ctx) error {
	currentUserID := utils.GetIDFromContext(c)
	if currentUserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, user needs to be logged in",
		})
	}

	idStr := utils.GetIDParam(c)
	if idStr == "" {
		return nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var board models.Board
	if err := config.DB.First(&board, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Board not found",
		})
	}

	// ✅ Ownership check
	if board.UserID != currentUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not own this board",
		})
	}

	if err := config.DB.Delete(&board).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete board",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Board deleted successfully",
	})
}

// ✅ GET BOARDS BY USER ID
func GetBoardByUserID(c *fiber.Ctx) error {
	idStr := utils.GetIDParam(c)
	if idStr == "" {
		return nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var boards []models.Board
	result := config.DB.Where("user_id = ?", uint(id)).Find(&boards)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch boards",
		})
	}

	return c.JSON(boards)
}