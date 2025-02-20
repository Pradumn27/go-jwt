package handlers

import (
	"context"
	"errors"
	"go-jwt/config"
	"go-jwt/models"
	"go-jwt/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

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
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No refresh token provided"})
	}

	// Parsing token with correct claims struct
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, errors.New("JWT_SECRET is not set")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Extracting claims
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Checking if token is expired
	exp, ok := (*claims)["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token expired"})
	}

	// Extracting user ID safely
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	// Generating new access token
	newAccessToken, _ := utils.GenerateToken(userID, time.Minute*15)

	// Setting the new access token in the response cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"message": "Token refreshed"})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "Tokens revoked"})
}