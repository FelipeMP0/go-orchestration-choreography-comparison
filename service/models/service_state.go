package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceState - Represents the current service state.
type ServiceState struct {
	ID        primitive.ObjectID `bson:"_id"`
	State     string             `bson:"state"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
