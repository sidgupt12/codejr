package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sidgupt12/codejr/models"
	"github.com/sidgupt12/codejr/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	log.Println("connected to db")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	app := fiber.New()
	routes.SetupRoutes(app, db)

	app.Listen(":" + os.Getenv("PORT"))

}
