package utils

import (
	"context"
	"go-jwt/config"
	"go-jwt/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
		"jti":     primitive.NewObjectID().Hex(), // Unique identifier for this token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func IsTokenRevoked(jti string) (bool, error) {
	var result models.BlacklistedToken
	err := config.DB.Collection("token_blacklist").FindOne(context.TODO(), bson.M{"jti": jti}).Decode(&result)
	return err == nil, err
}
