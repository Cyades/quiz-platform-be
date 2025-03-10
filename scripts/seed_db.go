package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Get MongoDB connection details
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	if mongoURI == "" {
		log.Fatal("MongoDB URI is not set. Please set MONGODB_URI environment variable")
	}
	if dbName == "" {
		dbName = "quiz_platform"
		log.Println("DB_NAME not specified, using default:", dbName)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}
	fmt.Println("Connected to MongoDB!")

	// Get the database and collection
	db := client.Database(dbName)
	collection := db.Collection("tryouts")

	// Drop the collection to start fresh
	if err := collection.Drop(ctx); err != nil {
		log.Fatal("Failed to drop collection: ", err)
	}
	fmt.Println("Dropped existing tryouts collection")

	// Create dummy tryout data
	dummyTryouts := []interface{}{
		bson.M{
			"title":       "Basic Mathematics Quiz",
			"description": "Test your basic math skills with this quiz covering arithmetic, algebra, and geometry concepts suitable for high school students.",
			"category":    "Mathematics",
			"duration":    30,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "English Grammar Challenge",
			"description": "Improve your grammar skills with this comprehensive quiz covering punctuation, sentence structure, and common English usage errors.",
			"category":    "Language",
			"duration":    45,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Science Fundamentals",
			"description": "Explore basic scientific concepts across physics, chemistry, and biology. Perfect for students preparing for general science tests.",
			"category":    "Science",
			"duration":    60,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "World History Overview",
			"description": "Test your knowledge of major historical events, civilizations, and influential figures throughout world history.",
			"category":    "History",
			"duration":    40,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Computer Science Basics",
			"description": "A quiz covering fundamental computer science concepts including algorithms, data structures, and basic programming principles.",
			"category":    "Computer Science",
			"duration":    50,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Geography Challenge",
			"description": "Test your knowledge of world geography, including countries, capitals, major landmarks, and geographical features.",
			"category":    "Geography",
			"duration":    35,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Physics Problem Solving",
			"description": "Challenge yourself with physics problem-solving scenarios covering mechanics, thermodynamics, and electromagnetism.",
			"category":    "Science",
			"duration":    55,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Literature Classics Quiz",
			"description": "Test your knowledge of classic literature, famous authors, literary movements, and iconic quotes from renowned works.",
			"category":    "Literature",
			"duration":    40,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Web Development Fundamentals",
			"description": "A quiz covering HTML, CSS, JavaScript, and basic web development concepts for beginners.",
			"category":    "Computer Science",
			"duration":    45,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		bson.M{
			"title":       "Economic Principles Test",
			"description": "Evaluate your understanding of basic economic concepts, theories, and real-world applications of economic principles.",
			"category":    "Economics",
			"duration":    50,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	// Insert the dummy data
	res, err := collection.InsertMany(ctx, dummyTryouts)
	if err != nil {
		log.Fatal("Failed to insert dummy data: ", err)
	}

	fmt.Printf("Successfully inserted %d dummy tryouts with IDs:\n", len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		fmt.Printf("%d: %s\n", i+1, id)
	}
}
