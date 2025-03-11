package controllers

import (
	"context"
	"log"
	"net/http"
	"quiz-platform/config"
	"quiz-platform/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const questionCollection = "questions"

// GetQuestionsByTryoutID returns all questions for a specific tryout
func GetQuestionsByTryoutID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tryoutID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(tryoutID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID format"})
		return
	}

	collection := config.GetCollection(questionCollection)
	findOptions := options.Find()
	findOptions.SetMaxTime(15 * time.Second)

	cursor, err := collection.Find(ctx, bson.M{"tryoutId": objectID}, findOptions)
	if err != nil {
		log.Printf("Error fetching questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions: " + err.Error()})
		return
	}

	var questions []models.Question
	if err = cursor.All(ctx, &questions); err != nil {
		log.Printf("Error decoding questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode questions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// GetQuestionByID returns a specific question by its ID
func GetQuestionByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tryoutID := c.Param("id")
	tryoutObjectID, err := primitive.ObjectIDFromHex(tryoutID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID format"})
		return
	}

	questionID := c.Param("questionId")
	questionObjectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	collection := config.GetCollection(questionCollection)
	findOneOptions := options.FindOne()
	findOneOptions.SetMaxTime(15 * time.Second)

	var question models.Question
	err = collection.FindOne(
		ctx,
		bson.M{
			"_id":      questionObjectID,
			"tryoutId": tryoutObjectID,
		},
		findOneOptions,
	).Decode(&question)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		log.Printf("Error fetching question %s: %v", questionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// CreateQuestion creates a new question for a tryout
func CreateQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tryoutID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(tryoutID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID format"})
		return
	}

	// Check if tryout exists and has no submissions
	tryoutCollection := config.GetCollection(tryoutCollection)
	var tryout models.Tryout
	err = tryoutCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tryout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryout: " + err.Error()})
		return
	}

	// Check if tryout has submissions
	if tryout.HasSubmission {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add questions to a tryout that has submissions"})
		return
	}

	var input models.QuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	now := time.Now()
	newQuestion := models.Question{
		TryoutID:  objectID,
		Text:      input.Text,
		IsTrue:    input.IsTrue,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection := config.GetCollection(questionCollection)
	insertOptions := options.InsertOne()
	result, err := collection.InsertOne(ctx, newQuestion, insertOptions)

	if err != nil {
		log.Printf("Error creating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question: " + err.Error()})
		return
	}

	newQuestion.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, newQuestion)
}

// UpdateQuestion updates an existing question
func UpdateQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	questionID := c.Param("questionId")
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	// Check if question exists and get tryout ID
	collection := config.GetCollection(questionCollection)
	var existingQuestion models.Question
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingQuestion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question: " + err.Error()})
		return
	}

	// Check if tryout has submissions
	tryoutCollection := config.GetCollection(tryoutCollection)
	var tryout models.Tryout
	err = tryoutCollection.FindOne(ctx, bson.M{"_id": existingQuestion.TryoutID}).Decode(&tryout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryout: " + err.Error()})
		return
	}

	if tryout.HasSubmission {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify questions of a tryout that has submissions"})
		return
	}

	var input models.QuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"text":      input.Text,
			"isTrue":    input.IsTrue,
			"updatedAt": time.Now(),
		},
	}

	updateOptions := options.Update()
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
		updateOptions,
	)

	if err != nil {
		log.Printf("Error updating question %s: %v", questionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Get updated question
	var updatedQuestion models.Question
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedQuestion)
	if err != nil {
		log.Printf("Error fetching updated question %s: %v", questionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Question updated but failed to retrieve updated data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedQuestion)
}

// DeleteQuestion deletes a question
func DeleteQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	questionID := c.Param("questionId")
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	// Check if question exists and get tryout ID
	collection := config.GetCollection(questionCollection)
	var existingQuestion models.Question
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingQuestion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question: " + err.Error()})
		return
	}

	// Check if tryout has submissions
	tryoutCollection := config.GetCollection(tryoutCollection)
	var tryout models.Tryout
	err = tryoutCollection.FindOne(ctx, bson.M{"_id": existingQuestion.TryoutID}).Decode(&tryout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryout: " + err.Error()})
		return
	}

	if tryout.HasSubmission {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete questions of a tryout that has submissions"})
		return
	}

	deleteOptions := options.Delete()
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID}, deleteOptions)
	if err != nil {
		log.Printf("Error deleting question %s: %v", questionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
