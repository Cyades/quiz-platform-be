package config

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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database connection variables
var (
	Client *mongo.Client
	DB     *mongo.Database
)

// ConnectDB initializes the MongoDB connection
func ConnectDB() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will use environment variables")
	}

	// Create MongoDB connection URI from env variables
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	if mongoURI == "" {
		log.Fatal("MongoDB URI is not set. Please set MONGODB_URI environment variable")
	}
	if dbName == "" {
		dbName = "quiz_platform"
		log.Println("DB_NAME not specified, using default:", dbName)
	}

	// Set client options with increased timeout for Atlas
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(20 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB Atlas...")
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}

	// Check the connection
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()

	err = Client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	DB = Client.Database(dbName)
	fmt.Println("Successfully connected to MongoDB Atlas!")
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			fmt.Println("Connection to MongoDB closed successfully.")
		}
	}
}

// SeedDummyData inserts sample data into the database if no data exists
func SeedDummyData() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := GetCollection("tryouts")

	// Check if collection is empty
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error checking collection count: %v", err)
		return
	}

	// Skip seeding if data already exists
	if count > 0 {
		fmt.Println("Dummy data already exists, skipping seed.")
		return
	}

	fmt.Println("Database is empty. Seeding with dummy data...")

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
	}

	// Insert the dummy data
	insertOptions := options.InsertMany()
	res, err := collection.InsertMany(ctx, dummyTryouts, insertOptions)
	if err != nil {
		log.Printf("Error seeding dummy data: %v", err)
		return
	}

	fmt.Printf("Successfully seeded %d dummy tryouts\n", len(res.InsertedIDs))
}
