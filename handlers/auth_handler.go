package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sidgupt12/codejr/models"
	"github.com/sidgupt12/codejr/utils"
	"gorm.io/gorm"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Parse Body
		var body registerRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		//validation
		if body.Name == "" || body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All fields are neccesary",
			})
		}

		//If user exists
		var existing models.User
		if err := db.Where("email = ?", body.Email).First(&existing).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User already exist, please login",
			})
		}

		//hash password
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error setting password",
			})
		}

		//create user
		user := models.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: hashedPassword,
		}

		//return success
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	}
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var body loginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "error logging in",
			})
		}

		//validation
		if body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All fields are neccesary",
			})
		}

		//find user
		var user models.User
		if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		// compare password
		if err := utils.CheckPassword(body.Password, user.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		token, err := utils.GenerateJWT(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "some error occured please try again later",
			})
		}

		return c.JSON(fiber.Map{
			"token": token,
		})

	}
}
