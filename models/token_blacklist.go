package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlacklistedToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	JTI       string             `bson:"jti"`
	ExpiresAt time.Time          `bson:"expires_at"`
}
