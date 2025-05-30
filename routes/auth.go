package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sidgupt12/codejr/handlers"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/register", handlers.Register(db))
	app.Post("/login", handlers.Login(db))
}
