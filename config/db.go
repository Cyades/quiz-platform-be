package config

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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get collections
	tryoutCollection := GetCollection("tryouts")
	questionCollection := GetCollection("questions")

	// Check if tryout collection is empty
	count, err := tryoutCollection.CountDocuments(ctx, bson.M{})
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

	// Create dummy tryout data with more realistic details
	dummyTryouts := []interface{}{
		bson.M{
			"title":         "Basic Mathematics Quiz",
			"description":   "Test your basic math skills with this quiz covering arithmetic, algebra, and geometry concepts suitable for high school students. Topics include equation solving, basic geometry theorems, and number properties.",
			"category":      "Mathematics",
			"duration":      30,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-10 * 24 * time.Hour), // 10 days ago
			"updatedAt":     time.Now().Add(-10 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Advanced Calculus Challenge",
			"description":   "Challenge yourself with complex calculus problems including limits, derivatives, integrals, and series. This tryout is designed for college-level mathematics students looking to test their understanding of advanced concepts.",
			"category":      "Mathematics",
			"duration":      60,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-15 * 24 * time.Hour), // 15 days ago
			"updatedAt":     time.Now().Add(-15 * 24 * time.Hour),
		},
		bson.M{
			"title":         "English Grammar Challenge",
			"description":   "Improve your grammar skills with this comprehensive quiz covering punctuation, sentence structure, verb tenses, and common English usage errors. Perfect for non-native speakers and language enthusiasts alike.",
			"category":      "Language",
			"duration":      45,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-8 * 24 * time.Hour), // 8 days ago
			"updatedAt":     time.Now().Add(-8 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Science Fundamentals",
			"description":   "Explore basic scientific concepts across physics, chemistry, and biology. Perfect for students preparing for general science tests. This quiz covers scientific method, basic laws of physics, periodic table concepts, and biological systems.",
			"category":      "Science",
			"duration":      60,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-5 * 24 * time.Hour), // 5 days ago
			"updatedAt":     time.Now().Add(-5 * 24 * time.Hour),
		},
		bson.M{
			"title":         "World History Overview",
			"description":   "Test your knowledge of major historical events, civilizations, and influential figures throughout world history. From ancient civilizations to modern geopolitics, this comprehensive quiz covers key moments that shaped our world.",
			"category":      "History",
			"duration":      40,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-12 * 24 * time.Hour), // 12 days ago
			"updatedAt":     time.Now().Add(-12 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Computer Science Basics",
			"description":   "A quiz covering fundamental computer science concepts including algorithms, data structures, and basic programming principles. Ideal for students beginning their journey into computer science or programmers wanting to review core concepts.",
			"category":      "Computer Science",
			"duration":      50,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-3 * 24 * time.Hour), // 3 days ago
			"updatedAt":     time.Now().Add(-3 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Geography Challenge",
			"description":   "Test your knowledge of world geography, including countries, capitals, major landmarks, and geographical features. This quiz will take you around the globe, from the highest peaks to the deepest oceans.",
			"category":      "Geography",
			"duration":      35,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
			"updatedAt":     time.Now().Add(-2 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Physics Problem Solving",
			"description":   "Challenge yourself with physics problem-solving scenarios covering mechanics, thermodynamics, and electromagnetism. This advanced quiz requires application of physics principles to solve complex, real-world problems.",
			"category":      "Science",
			"duration":      55,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-7 * 24 * time.Hour), // 7 days ago
			"updatedAt":     time.Now().Add(-7 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Literature Classics Quiz",
			"description":   "Test your knowledge of classic literature, famous authors, literary movements, and iconic quotes from renowned works. From Shakespeare to Tolstoy, this quiz covers literary masterpieces from around the world.",
			"category":      "Literature",
			"duration":      40,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-9 * 24 * time.Hour), // 9 days ago
			"updatedAt":     time.Now().Add(-9 * 24 * time.Hour),
		},
		bson.M{
			"title":         "Web Development Fundamentals",
			"description":   "A quiz covering HTML, CSS, JavaScript, and basic web development concepts for beginners. Test your understanding of responsive design, DOM manipulation, and basic front-end development principles.",
			"category":      "Computer Science",
			"duration":      45,
			"hasSubmission": false,
			"createdAt":     time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
			"updatedAt":     time.Now().Add(-1 * 24 * time.Hour),
		},
	}

	// Insert tryout data first
	tryoutResults, err := tryoutCollection.InsertMany(ctx, dummyTryouts)
	if err != nil {
		log.Printf("Error seeding tryout data: %v", err)
		return
	}
	fmt.Printf("Successfully seeded %d dummy tryouts\n", len(tryoutResults.InsertedIDs))

	// Now create questions for each tryout
	var questions []interface{}

	// Math questions for "Basic Mathematics Quiz"
	math1ID := tryoutResults.InsertedIDs[0].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  math1ID,
			"text":      "The square root of 144 is 12.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math1ID,
			"text":      "In a right-angled triangle, the square of the hypotenuse equals the sum of the squares of the other two sides.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math1ID,
			"text":      "The formula for the area of a circle is πr.",
			"isTrue":    false, // It's πr²
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math1ID,
			"text":      "The sum of all angles in a triangle is 180 degrees.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math1ID,
			"text":      "The value of π (pi) is exactly 22/7.",
			"isTrue":    false, // It's an irrational number, 22/7 is an approximation
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// Advanced Calculus questions
	math2ID := tryoutResults.InsertedIDs[1].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  math2ID,
			"text":      "The derivative of e^x is e^x.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math2ID,
			"text":      "The integral of 1/x is ln|x| + C.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math2ID,
			"text":      "For any continuous function f(x), the derivative of the integral of f(x) from a to x with respect to x is equal to f(x).",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  math2ID,
			"text":      "L'Hôpital's rule can be applied to any indeterminate form.",
			"isTrue":    false, // Only applicable to 0/0 and ∞/∞ forms directly
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// English Grammar questions
	englishID := tryoutResults.InsertedIDs[2].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  englishID,
			"text":      "In English, the subject always comes before the verb in a sentence.",
			"isTrue":    false, // Not in questions or certain literary constructions
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  englishID,
			"text":      "'i' comes before 'e' except after 'c' is a grammar rule that has no exceptions.",
			"isTrue":    false, // Many exceptions like "weird", "science", "efficient"
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  englishID,
			"text":      "A semicolon can be used to join two independent clauses.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  englishID,
			"text":      "The past participle of 'go' is 'went'.",
			"isTrue":    false, // It's "gone", "went" is past tense
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// Science questions
	scienceID := tryoutResults.InsertedIDs[3].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  scienceID,
			"text":      "Mitochondria are known as the powerhouse of the cell.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  scienceID,
			"text":      "According to Newton's First Law, an object will remain at rest or in uniform motion unless acted upon by an external force.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  scienceID,
			"text":      "Water's chemical formula is H2O2.",
			"isTrue":    false, // It's H2O, H2O2 is hydrogen peroxide
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  scienceID,
			"text":      "DNA is a double helix structure.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  scienceID,
			"text":      "Sound travels faster in air than in water.",
			"isTrue":    false, // Sound travels faster in denser mediums like water
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// History questions
	historyID := tryoutResults.InsertedIDs[4].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  historyID,
			"text":      "The American Declaration of Independence was signed in 1776.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  historyID,
			"text":      "The Berlin Wall fell in 1989.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  historyID,
			"text":      "World War II ended in 1950.",
			"isTrue":    false, // It ended in 1945
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  historyID,
			"text":      "The Ancient Roman Empire was centered around Greece.",
			"isTrue":    false, // It was centered around Rome, Italy
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// Computer Science questions
	csID := tryoutResults.InsertedIDs[5].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  csID,
			"text":      "In binary, the decimal number 10 is represented as 1010.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  csID,
			"text":      "HTML is a programming language.",
			"isTrue":    false, // It's a markup language
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  csID,
			"text":      "An array index typically starts at 1 in most programming languages.",
			"isTrue":    false, // Most start at 0
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  csID,
			"text":      "The Big O notation O(n²) represents a quadratic time complexity.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  csID,
			"text":      "DNS stands for Domain Name System.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// Web Dev questions
	webdevID := tryoutResults.InsertedIDs[9].(primitive.ObjectID)
	questions = append(questions, []interface{}{
		bson.M{
			"tryoutId":  webdevID,
			"text":      "CSS stands for Cascading Style Sheets.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  webdevID,
			"text":      "JavaScript can directly modify database records without a backend server.",
			"isTrue":    false, // Client-side JS cannot directly access databases
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  webdevID,
			"text":      "The box model in CSS consists of margin, border, padding, and content.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  webdevID,
			"text":      "HTTP status code 404 means 'Server Error'.",
			"isTrue":    false, // 404 is "Not Found", 500 is "Server Error"
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"tryoutId":  webdevID,
			"text":      "In responsive design, the 'viewport' meta tag helps to ensure proper display on mobile devices.",
			"isTrue":    true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}...)

	// Insert question data
	questionResult, err := questionCollection.InsertMany(ctx, questions)
	if err != nil {
		log.Printf("Error seeding question data: %v", err)
		return
	}
	fmt.Printf("Successfully seeded %d dummy questions\n", len(questionResult.InsertedIDs))
}
