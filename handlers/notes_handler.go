package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sidgupt12/codejr/models"
	"gorm.io/gorm"
)

func CreateNote(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)

		var note models.Note
		if err := c.BodyParser(&note); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		if note.Title == "" || note.Content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Both fields are mandatory",
			})
		}

		note.UserId = userID

		if err := db.Create(&note).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create note",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(note)

	}
}

func GetNotes(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)

		pageStr := c.Query("page", "1")
		limitStr := c.Query("limit", "10")
		search := c.Query("search", "")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var notes []models.Note
		query := db.Where("user_id = ?", userID)

		if search != "" {
			likeSearch := "%" + search + "%"
			query = query.Where("title LIKE ?", likeSearch)
		}

		if err := query.
			Limit(limit).
			Offset(offset).
			Find(&notes).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch notes",
			})
		}

		return c.JSON(notes)
	}
}

func UpdateNote(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)
		noteID := c.Params("id")

		var note models.Note
		if err := db.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}

		var updateData struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		if updateData.Title == "" || updateData.Content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Both fields are mandatory",
			})
		}

		note.Title = updateData.Title
		note.Content = updateData.Content

		if err := db.Save(&note).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update note",
			})
		}

		return c.JSON(note)
	}
}

func DeleteNote(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)
		noteID := c.Params("id")

		var note models.Note
		if err := db.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}

		if err := db.Delete(&note).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete note",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Note deleted successfully",
		})
	}
}
