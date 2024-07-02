package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"wagobot.com/controllers"
	"wagobot.com/db"
	"wagobot.com/middleware"
	"wagobot.com/router"
)

func main() {
	// Load values from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve values from environment variables
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")

	// Configure database logging
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Initialize SQL store
	storeContainer, err := sqlstore.New("sqlite3", dbPath, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set storeContainer to the controller package variable
	controllers.SetStoreContainer(storeContainer)

	// Initialize the database
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.CloseDB()

	// Setup router
	r := router.SetupRouter()
	apiRouter := router.SetupRouter()
	r.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Enable CORS for development
	corsHandler := middleware.SetupCORS(r)

	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
