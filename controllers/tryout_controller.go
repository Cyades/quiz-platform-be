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

const tryoutCollection = "tryouts"

// GetAllTryouts returns all tryouts
func GetAllTryouts(c *gin.Context) {
	// Increase timeout for MongoDB Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := config.GetCollection(tryoutCollection)

	// Set options to handle potential timeout issues
	findOptions := options.Find()
	findOptions.SetMaxTime(15 * time.Second)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)

	if err != nil {
		log.Printf("Error fetching tryouts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryouts: " + err.Error()})
		return
	}

	var tryouts []models.Tryout
	if err = cursor.All(ctx, &tryouts); err != nil {
		log.Printf("Error decoding tryouts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tryouts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tryouts)
}

// GetTryout returns a specific tryout by ID
func GetTryout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var tryout models.Tryout
	collection := config.GetCollection(tryoutCollection)

	// Set options for FindOne operation
	findOneOptions := options.FindOne()
	findOneOptions.SetMaxTime(15 * time.Second)

	err = collection.FindOne(ctx, bson.M{"_id": objectID}, findOneOptions).Decode(&tryout)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
			return
		}
		log.Printf("Error fetching tryout %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tryout)
}

// CreateTryout creates a new tryout
func CreateTryout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var input models.TryoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	now := time.Now()
	newTryout := models.Tryout{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		Duration:    input.Duration,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	collection := config.GetCollection(tryoutCollection)
	insertOptions := options.InsertOne()
	result, err := collection.InsertOne(ctx, newTryout, insertOptions)

	if err != nil {
		log.Printf("Error creating tryout: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tryout: " + err.Error()})
		return
	}

	newTryout.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, newTryout)
}

// UpdateTryout updates an existing tryout
func UpdateTryout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input models.TryoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       input.Title,
			"description": input.Description,
			"category":    input.Category,
			"duration":    input.Duration,
			"updatedAt":   time.Now(),
		},
	}

	collection := config.GetCollection(tryoutCollection)
	updateOptions := options.Update()
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
		updateOptions,
	)

	if err != nil {
		log.Printf("Error updating tryout %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tryout: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}

	// Get updated tryout
	findOneOptions := options.FindOne()
	findOneOptions.SetMaxTime(15 * time.Second)

	var updatedTryout models.Tryout
	err = collection.FindOne(ctx, bson.M{"_id": objectID}, findOneOptions).Decode(&updatedTryout)
	if err != nil {
		log.Printf("Error fetching updated tryout %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tryout updated but failed to retrieve updated data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTryout)
}

// DeleteTryout deletes a tryout
func DeleteTryout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.GetCollection(tryoutCollection)
	deleteOptions := options.Delete()
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID}, deleteOptions)
	if err != nil {
		log.Printf("Error deleting tryout %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tryout: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tryout deleted successfully"})
}

// GetTryoutOptions returns all possible categories for filtering (helper function)
func GetTryoutOptions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection(tryoutCollection)

	// Get unique categories
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$category"}},
		{"$project": bson.M{"category": "$_id", "_id": 0}},
	}
	categoryCursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var categories []bson.M
	if err = categoryCursor.All(ctx, &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
