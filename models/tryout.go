package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tryout represents a quiz/tryout in the platform
type Tryout struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	Category      string             `json:"category" bson:"category"`
	Duration      int                `json:"duration" bson:"duration"` // in minutes
	HasSubmission bool               `json:"hasSubmission" bson:"hasSubmission"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// TryoutInput is used for creating or updating a tryout
type TryoutInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Duration    int    `json:"duration" binding:"required,min=1"`
}
