package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sidgupt12/codejr/models"
	"github.com/sidgupt12/codejr/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


// app was starting before the database was ready, so we need to retry
func connectWithRetry(dsn string, retries int, delay time.Duration) (*gorm.DB, error) {
  var db *gorm.DB
  var err error
  for i := 0; i < retries; i++ {
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err == nil {
      return db, nil
    }
    log.Printf("DB connect failed: %v. Retrying in %v...", err, delay)
    time.Sleep(delay)
  }
  return nil, err
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URL")

	db, err := connectWithRetry(dsn, 10, 3*time.Second)  
	if err != nil {
		log.Fatal("Error connecting to database after retries:", err)
	}
	log.Println("connected to db")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	err = db.AutoMigrate(&models.Note{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	app := fiber.New()
	routes.SetupRoutes(app, db)

	app.Listen("0.0.0.0:" + os.Getenv("PORT"))

}
