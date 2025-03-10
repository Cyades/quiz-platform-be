# Quiz Platform Backend

This is the backend component of the Quiz Platform application, built with Go and MongoDB.

## Features

- RESTful API for managing quiz tryouts
- MongoDB integration for data persistence
- Complete CRUD operations for tryouts
- Automatic database seeding with sample data

## Prerequisites

- Go (v1.16 or later)
- MongoDB (v4.4 or later) or MongoDB Atlas account

## Setup and Installation

### 1. Install Dependencies

```bash
# Install all required Go dependencies
go mod download
```

### 2. Configure Environment

The application uses environment variables for configuration. A `.env` file is provided:

```
PORT=8080
MONGODB_URI=mongodb+srv://malvinscafi:4VRgd29BBJK24xNo@cluster0.chomd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
DB_NAME=quiz_platform
```

If you prefer using a local MongoDB installation, modify the file:

```
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DB_NAME=quiz_platform
```

### 3. Start MongoDB (if using local instance)

If you're using a local MongoDB installation, make sure it's running:

```bash
# MongoDB typically runs on port 27017
mongod --dbpath /path/to/data/db
```

### 4. Run the Application

```bash
go run main.go
```

The server will:
1. Connect to MongoDB
2. Seed the database with sample data if empty
3. Start the HTTP server on the configured port (default: 8080)


## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /api/v1/tryouts | Get all tryouts |
| GET    | /api/v1/tryouts/:id | Get a specific tryout by ID |
| POST   | /api/v1/tryouts | Create a new tryout |
| PUT    | /api/v1/tryouts/:id | Update an existing tryout |
| DELETE | /api/v1/tryouts/:id | Delete a tryout |
| GET    | /api/v1/tryouts/filter/options | Get filtering options |


## Seeding Data

To reset and repopulate the database with sample data:

```bash
cd scripts
go run seed_db.go
```

## Project Structure

```
backend/
├── config/         # Database configuration
│   └── db.go       # MongoDB connection setup
├── controllers/    # API controllers
│   └── tryout_controller.go  # Tryout endpoints
├── models/         # Data models
│   └── tryout.go   # Tryout data structure
├── routes/         # API routes
│   └── routes.go   # Route definitions
├── scripts/        # Utility scripts
│   └── seed_db.go  # Database seeding
├── .env            # Environment variables
├── go.mod          # Go module file
├── go.sum          # Go dependencies checksums
└── main.go         # Application entry point
```