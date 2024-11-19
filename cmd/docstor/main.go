package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"docstor-be/api/router"
	"docstor-be/internal/db"
	"docstor-be/internal/document"
	"docstor-be/internal/user"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Prefix for database-related environment variables
	prefix := "DB_"

	// Initialize connection string
	connStr := ""

	// Loop through all environment variables
	for _, env := range os.Environ() {
		// Split the environment variable into key and value
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		value := pair[1]

		// Check if the key has the DB_ prefix
		if strings.HasPrefix(key, prefix) {
			// Format the key for connection string (e.g., "DB_HOST" to "host")
			connKey := strings.ToLower(strings.TrimPrefix(key, prefix))
			connStr += fmt.Sprintf("%s=%s ", connKey, value)
		}
	}

	// Trim any trailing spaces
	connStr = strings.TrimSpace(connStr)

	fmt.Println("Connection String:", connStr)

	// Set up the database connection
	dbConn, err := db.ConnectToDatabase(connStr)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer dbConn.Close()

	// create firebase app instance
	opt := option.WithCredentialsFile("config/service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	// create firebase auth instance
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	// run database migration
	dbConn.AutoMigrate(&user.User{})

	// setup service
	userService := &user.UserService{
		DB:       dbConn,
		FireAuth: authClient,
	}
	documentService := document.NewDocumentService()

	// setup handler
	userHandler := user.NewUserHandler(userService)
	documentHandler := document.NewDocumentHandler(documentService)

	// setup http router
	r := router.NewRouter(userHandler, documentHandler)

	r.Run(":8080")
}
