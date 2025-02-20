package handlers

import (
	"context"
	"go-jwt/config"
	"go-jwt/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDetails(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userObjectID, userIdErr := primitive.ObjectIDFromHex(userID)
	if userIdErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authorized"})
	}

	return c.JSON(fiber.Map{ "user": user})
}