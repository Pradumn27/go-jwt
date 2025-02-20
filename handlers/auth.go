package handlers

import (
	"context"
	"go-jwt/config"
	"go-jwt/models"
	"go-jwt/utils"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPassword, hashingErr := utils.HashPassword(user.Password)
	if hashingErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Password hashing failed"})
	}
	user.Password = hashedPassword

	_, err := config.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User creation failed"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created"})
}

func SignIn(c *fiber.Ctx) error {
	var user models.User
	input := new(models.User)
	
	parsingErr := c.BodyParser(input)
	if parsingErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, _ := utils.GenerateToken(user.ID.Hex())
	return c.JSON(fiber.Map{"token": token})
}

func RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	tokenString := tokenParts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	jti, _ := claims["jti"].(string)
	exp, _ := claims["exp"].(float64)

	revoked, _ := utils.IsTokenRevoked(jti)
	if revoked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has been revoked"})
	}

	userID, exists := claims["user_id"].(string)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token payload"})
	}

	// Revoke the previous token
	_, err = config.DB.Collection("token_blacklist").InsertOne(context.TODO(), models.BlacklistedToken{
		JTI:       jti,
		ExpiresAt: time.Unix(int64(exp), 0),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to revoke previous token"})
	}

	newToken, err := utils.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": newToken})
}

func RevokeToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	tokenString := tokenParts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	jti, _ := claims["jti"].(string)
	exp, _ := claims["exp"].(float64)

	// Store the JTI in the blacklist collection
	_, err = config.DB.Collection("token_blacklist").InsertOne(context.TODO(), models.BlacklistedToken{
		JTI:       jti,
		ExpiresAt: time.Unix(int64(exp), 0),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to revoke token"})
	}

	return c.JSON(fiber.Map{"message": "Token revoked successfully"})
}
