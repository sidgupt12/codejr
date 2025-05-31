package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sidgupt12/codejr/models"
	"gorm.io/gorm"
)

type NoteResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

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

		var total int64
		query.Model(&models.Note{}).Count(&total)

		if err := query.
			Limit(limit).
			Offset(offset).
			Find(&notes).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch notes",
			})
		}

		// Convert to response format
		var response []NoteResponse
		for _, note := range notes {
			response = append(response, NoteResponse{
				ID:      note.ID,
				Title:   note.Title,
				Content: note.Content,
			})
		}

		return c.JSON(fiber.Map{
			"data": response,
			"meta": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": int((total + int64(limit) - 1) / int64(limit)),
			},
		})
	}
}

func GetNoteById(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)
		noteID := c.Params("id")

		var note models.Note
		if err := db.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Note not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch note",
			})
		}
		return c.JSON(fiber.Map{
			"id":      note.ID,
			"title":   note.Title,
			"content": note.Content,
		})
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
