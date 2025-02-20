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

/*

	SignUp: 
	curl -X POST http://localhost:3000/auth/signup \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "securepassword"
         }' \

	SignIn: 
	curl -X POST http://localhost:3000/auth/signin \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "securepassword"
         }' \

	Refresh: 
	curl -X GET http://localhost:3000/auth/refresh \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNTgwNzQsImp0aSI6IjY3YjcyNmQyMmVkMDcyYWZjNWRiMjQyMyIsInVzZXJfaWQiOiI2N2I3MjY5ZTJlZDA3MmFmYzVkYjI0MjIifQ.jT4IXbyGcZPPZUS_zjc0_M1jlwHrFCFIHiLjQbhxoU0"


	 UserDetails
	 curl -X GET http://localhost:3000/user/details \
	    -H "Content-Type: application/json"  \
		-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNTgwNzQsImp0aSI6IjY3YjcyNmQyMmVkMDcyYWZjNWRiMjQyMyIsInVzZXJfaWQiOiI2N2I3MjY5ZTJlZDA3MmFmYzVkYjI0MjIifQ.jT4IXbyGcZPPZUS_zjc0_M1jlwHrFCFIHiLjQbhxoU0"

	RevokeToken: 
	curl -X POST http://localhost:3000/auth/revoke \
     	-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNTgwNzQsImp0aSI6IjY3YjcyNmQyMmVkMDcyYWZjNWRiMjQyMyIsInVzZXJfaWQiOiI2N2I3MjY5ZTJlZDA3MmFmYzVkYjI0MjIifQ.jT4IXbyGcZPPZUS_zjc0_M1jlwHrFCFIHiLjQbhxoU0"

*/