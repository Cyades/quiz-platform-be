package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents a question in a tryout
type Question struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TryoutID  primitive.ObjectID `json:"tryoutId" bson:"tryoutId"`
	Text      string             `json:"text" bson:"text"`
	IsTrue    bool               `json:"isTrue" bson:"isTrue"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// QuestionInput is used for creating or updating a question
type QuestionInput struct {
	Text   string `json:"text" binding:"required"`
	IsTrue bool   `json:"isTrue"`
}
