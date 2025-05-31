package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sidgupt12/codejr/handlers"
	"github.com/sidgupt12/codejr/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// Auth Routes
	app.Post("/register", handlers.Register(db))
	app.Post("/login", handlers.Login(db))

	// Note Routes

	// protecting routes by forming groups
	notes := app.Group("/note", middleware.Protect())

	notes.Post("/", handlers.CreateNote(db))
	notes.Get("/", handlers.GetNotes(db))
	notes.Put("/:id", handlers.UpdateNote(db))
	notes.Delete("/:id", handlers.DeleteNote(db))

}
