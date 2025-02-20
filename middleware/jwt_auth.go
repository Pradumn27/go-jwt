package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTAuthGuard(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}
		tokenString := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		userID, ok1 := claims["user_id"].(string)
		jti, ok2 := claims["jti"].(string)
		if !ok1 || !ok2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Check if the token's JTI is in the blacklist
		var result bson.M
		err = db.Collection("token_blacklist").FindOne(context.TODO(), bson.M{"jti": jti}).Decode(&result)
		if err == nil {
			// Token exists in blacklist, reject it
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token revoked"})
		}

		// Set user info in context for later use
		c.Locals("user_id", userID)
		c.Locals("jti", jti)

		return c.Next()
	}
}
