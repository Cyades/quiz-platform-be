package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Get the database and collections
	db := client.Database(dbName)
	tryoutCollection := db.Collection("tryouts")
	questionCollection := db.Collection("questions")

	// Drop the collections to start fresh
	if err := tryoutCollection.Drop(ctx); err != nil {
		log.Fatal("Failed to drop tryouts collection: ", err)
	}
	fmt.Println("Dropped existing tryouts collection")

	if err := questionCollection.Drop(ctx); err != nil {
		log.Fatal("Failed to drop questions collection: ", err)
	}
	fmt.Println("Dropped existing questions collection")

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
			"createdAt":   time.Now().Add(-24 * time.Hour),
			"updatedAt":   time.Now().Add(-24 * time.Hour),
		},
		bson.M{
			"title":       "Science Fundamentals",
			"description": "Explore basic scientific concepts across physics, chemistry, and biology. Perfect for students preparing for general science tests.",
			"category":    "Science",
			"duration":    60,
			"createdAt":   time.Now().Add(-48 * time.Hour),
			"updatedAt":   time.Now().Add(-48 * time.Hour),
		},
		bson.M{
			"title":       "World History Overview",
			"description": "Test your knowledge of major historical events, civilizations, and influential figures throughout world history.",
			"category":    "History",
			"duration":    40,
			"createdAt":   time.Now().Add(-72 * time.Hour),
			"updatedAt":   time.Now().Add(-72 * time.Hour),
		},
		bson.M{
			"title":       "Computer Science Basics",
			"description": "A quiz covering fundamental computer science concepts including algorithms, data structures, and basic programming principles.",
			"category":    "Computer Science",
			"duration":    50,
			"createdAt":   time.Now().Add(-96 * time.Hour),
			"updatedAt":   time.Now().Add(-96 * time.Hour),
		},
		bson.M{
			"title":       "Geography Challenge",
			"description": "Test your knowledge of world geography, including countries, capitals, major landmarks, and geographical features.",
			"category":    "Geography",
			"duration":    35,
			"createdAt":   time.Now().Add(-120 * time.Hour),
			"updatedAt":   time.Now().Add(-120 * time.Hour),
		},
		bson.M{
			"title":       "Physics Problem Solving",
			"description": "Challenge yourself with physics problem-solving scenarios covering mechanics, thermodynamics, and electromagnetism.",
			"category":    "Science",
			"duration":    55,
			"createdAt":   time.Now().Add(-144 * time.Hour),
			"updatedAt":   time.Now().Add(-144 * time.Hour),
		},
		bson.M{
			"title":       "Literature Classics Quiz",
			"description": "Test your knowledge of classic literature, famous authors, literary movements, and iconic quotes from renowned works.",
			"category":    "Literature",
			"duration":    40,
			"createdAt":   time.Now().Add(-168 * time.Hour),
			"updatedAt":   time.Now().Add(-168 * time.Hour),
		},
		bson.M{
			"title":       "Web Development Fundamentals",
			"description": "A quiz covering HTML, CSS, JavaScript, and basic web development concepts for beginners.",
			"category":    "Computer Science",
			"duration":    45,
			"createdAt":   time.Now().Add(-192 * time.Hour),
			"updatedAt":   time.Now().Add(-192 * time.Hour),
		},
		bson.M{
			"title":       "Economic Principles Test",
			"description": "Evaluate your understanding of basic economic concepts, theories, and real-world applications of economic principles.",
			"category":    "Economics",
			"duration":    50,
			"createdAt":   time.Now().Add(-216 * time.Hour),
			"updatedAt":   time.Now().Add(-216 * time.Hour),
		},
		// New additional tryouts
		bson.M{
			"title":       "Psychology Concepts Quiz",
			"description": "Test your understanding of fundamental psychology theories, famous experiments, and human behavior patterns.",
			"category":    "Psychology",
			"duration":    40,
			"createdAt":   time.Now().Add(-240 * time.Hour),
			"updatedAt":   time.Now().Add(-240 * time.Hour),
		},
		bson.M{
			"title":       "Mobile App Development Basics",
			"description": "Evaluate your knowledge of mobile development concepts, frameworks, and best practices for iOS and Android platforms.",
			"category":    "Computer Science",
			"duration":    60,
			"createdAt":   time.Now().Add(-264 * time.Hour),
			"updatedAt":   time.Now().Add(-264 * time.Hour),
		},
		bson.M{
			"title":       "Art History Through the Ages",
			"description": "Journey through different art periods, famous artists, and iconic works that have shaped the history of visual arts.",
			"category":    "Art",
			"duration":    45,
			"createdAt":   time.Now().Add(-288 * time.Hour),
			"updatedAt":   time.Now().Add(-288 * time.Hour),
		},
		bson.M{
			"title":       "Introduction to Philosophy",
			"description": "Explore major philosophical questions, influential thinkers, and fundamental concepts in Western and Eastern philosophy.",
			"category":    "Philosophy",
			"duration":    55,
			"createdAt":   time.Now().Add(-312 * time.Hour),
			"updatedAt":   time.Now().Add(-312 * time.Hour),
		},
		bson.M{
			"title":       "Renewable Energy Technologies",
			"description": "Test your knowledge of sustainable energy solutions, including solar, wind, hydro, and emerging green technologies.",
			"category":    "Science",
			"duration":    35,
			"createdAt":   time.Now().Add(-336 * time.Hour),
			"updatedAt":   time.Now().Add(-336 * time.Hour),
		},
	}

	// Insert the dummy tryout data
	tryoutRes, err := tryoutCollection.InsertMany(ctx, dummyTryouts)
	if err != nil {
		log.Fatal("Failed to insert dummy tryout data: ", err)
	}

	fmt.Printf("Successfully inserted %d dummy tryouts with IDs:\n", len(tryoutRes.InsertedIDs))

	// Create dummy questions for each tryout
	var questions []interface{}

	// Create different questions for each tryout
	mathQuestions := []bson.M{
		{
			"text":   "The Pythagorean theorem applies to all triangles.",
			"isTrue": false,
		},
		{
			"text":   "The sum of interior angles in a triangle is 180 degrees.",
			"isTrue": true,
		},
		{
			"text":   "A prime number is divisible only by 1 and itself.",
			"isTrue": true,
		},
		{
			"text":   "The square root of a negative number is always undefined.",
			"isTrue": false,
		},
		{
			"text":   "In calculus, the derivative of a constant is always zero.",
			"isTrue": true,
		},
	}

	grammarQuestions := []bson.M{
		{
			"text":   "In English grammar, a sentence must always contain a verb.",
			"isTrue": true,
		},
		{
			"text":   "The words 'their', 'there', and 'they're' are all pronounced the same.",
			"isTrue": true,
		},
		{
			"text":   "A semicolon is used to connect independent clauses that are closely related.",
			"isTrue": true,
		},
		{
			"text":   "Adjectives always come after the noun they modify in English.",
			"isTrue": false,
		},
		{
			"text":   "The passive voice should never be used in good writing.",
			"isTrue": false,
		},
	}

	scienceQuestions := []bson.M{
		{
			"text":   "Water's chemical formula is H₂O.",
			"isTrue": true,
		},
		{
			"text":   "Atoms are mostly empty space.",
			"isTrue": true,
		},
		{
			"text":   "The Earth's core is primarily composed of iron and nickel.",
			"isTrue": true,
		},
		{
			"text":   "Sound travels faster in air than in water.",
			"isTrue": false,
		},
		{
			"text":   "Photosynthesis is the process of converting light energy into chemical energy.",
			"isTrue": true,
		},
	}

	historyQuestions := []bson.M{
		{
			"text":   "The Roman Empire fell in 476 CE.",
			"isTrue": true,
		},
		{
			"text":   "The United States declared independence in 1776.",
			"isTrue": true,
		},
		{
			"text":   "The Magna Carta was signed in 1215.",
			"isTrue": true,
		},
		{
			"text":   "The French Revolution began in 1689.",
			"isTrue": false,
		},
		{
			"text":   "Christopher Columbus was the first European to reach North America.",
			"isTrue": false,
		},
	}

	csQuestions := []bson.M{
		{
			"text":   "HTML is a programming language.",
			"isTrue": false,
		},
		{
			"text":   "The binary system uses the digits 0 and 1.",
			"isTrue": true,
		},
		{
			"text":   "A byte consists of 8 bits.",
			"isTrue": true,
		},
		{
			"text":   "RAM is a type of non-volatile memory.",
			"isTrue": false,
		},
		{
			"text":   "All algorithms have a finite number of steps.",
			"isTrue": true,
		},
	}

	// Generic questions for the remaining tryouts
	genericQuestions := []bson.M{
		{
			"text":   "This statement is true.",
			"isTrue": true,
		},
		{
			"text":   "This statement is false.",
			"isTrue": false,
		},
		{
			"text":   "The quiz platform supports true/false questions.",
			"isTrue": true,
		},
		{
			"text":   "All questions must have exactly five possible answers.",
			"isTrue": false,
		},
		{
			"text":   "The quiz platform allows editing questions if there are no submissions.",
			"isTrue": true,
		},
	}

	// Map of question sets for the first 5 tryouts
	questionSets := map[int][]bson.M{
		0: mathQuestions,
		1: grammarQuestions,
		2: scienceQuestions,
		3: historyQuestions,
		4: csQuestions,
	}

	// Add questions for each tryout
	now := time.Now()
	for i, id := range tryoutRes.InsertedIDs {
		tryoutID := id.(primitive.ObjectID)

		// Choose the appropriate question set
		var questionSet []bson.M
		if i < 5 {
			questionSet = questionSets[i]
		} else {
			questionSet = genericQuestions
		}

		// Create questions for this tryout
		for _, q := range questionSet {
			questionTime := now.Add(time.Duration(-i) * 24 * time.Hour)

			question := bson.M{
				"tryoutId":  tryoutID,
				"text":      q["text"],
				"isTrue":    q["isTrue"],
				"createdAt": questionTime,
				"updatedAt": questionTime,
			}

			questions = append(questions, question)
		}

		fmt.Printf("%d: %s - Added %d questions\n", i+1, tryoutID.Hex(), len(questionSet))
	}

	// Insert all questions
	if len(questions) > 0 {
		questionRes, err := questionCollection.InsertMany(ctx, questions)
		if err != nil {
			log.Fatal("Failed to insert dummy question data: ", err)
		}
		fmt.Printf("Successfully inserted %d dummy questions\n", len(questionRes.InsertedIDs))
	}

	fmt.Println("Database seeded successfully!")
}
