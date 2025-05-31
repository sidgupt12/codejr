package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sidgupt12/codejr/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load env variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, make sure environment variables are set")
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Seed users
	users := []models.User{
		{Name: "Alice", Email: "alice@example.com", Password: hashPassword("password1")},
		{Name: "Bob", Email: "bob@example.com", Password: hashPassword("password2")},
	}

	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			log.Println("Failed to create user:", u.Email, err)
		}
	}

	notes := []models.Note{
		{UserId: 1, Title: "Alice Note 1", Content: "Content for Alice note 1"},
		{UserId: 1, Title: "Alice Note 2", Content: "Content for Alice note 2"},
		{UserId: 2, Title: "Bob Note 1", Content: "Content for Bob note 1"},
	}

	for _, n := range notes {
		if err := db.Create(&n).Error; err != nil {
			log.Println("Failed to create note:", n.Title, err)
		}
	}

	fmt.Println("Seeding done!")
}

func hashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return string(hash)
}
