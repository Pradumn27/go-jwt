package handlers

import (
	"context"
	"go-jwt/config"
	"go-jwt/models"
	"go-jwt/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	_, err := config.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User creation failed"})
	}

	accessToken, _ := utils.GenerateToken(user.ID.String(), time.Minute*15)
	refreshToken, _ := utils.GenerateToken(user.ID.String(), time.Hour*24*7)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created"})
}

func SignIn(c *fiber.Ctx) error {
	var user models.User
	input := new(models.User)
	_ = c.BodyParser(input)

	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	accessToken, _ := utils.GenerateToken(user.ID.String(), time.Minute*15)
	refreshToken, _ := utils.GenerateToken(user.ID.String(), time.Hour*24*7)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})


	return c.JSON(fiber.Map{"message": "Login successful"})
}

func RefreshToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Token Refresh successful"})
}
