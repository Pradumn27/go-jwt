package handlers

import (
	"context"
	"go-jwt/config"
	"go-jwt/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetDetails(c *fiber.Ctx) error {
	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": "27pradumn@gmail.com"}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{ "user": user})
}