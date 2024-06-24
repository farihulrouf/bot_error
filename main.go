package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"wagobot.com/controllers"
	"wagobot.com/router"
)

func main() {
	// Load values from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve values from environment variables
	//secretKey := os.Getenv("SECRET_KEY")
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH") // Get DB_PATH from .env

	// Configure database logging
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	index := strings.Index(dbPath, ":")
	if index == -1 {
		fmt.Println("Format dbPath tidak valid")
		return
	}

	// Ambil substring setelah titik dua dan sebelum tanda tanya (?)
	substring := dbPath[index+1:]
	endIndex := strings.Index(substring, "?")
	if endIndex == -1 {
		fmt.Println("Format dbPath tidak valid")
		return
	}

	fileName := substring[:endIndex]
	err = removeFile(fileName)
	if err != nil {
		log.Printf("Failed to remove SQLite database file: %v", err)
	}

	// Initialize SQL store
	storeContainer, err := sqlstore.New("sqlite3", dbPath, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set storeContainer to the controller package variable
	controllers.SetStoreContainer(storeContainer)

	// Setup router
	r := router.SetupRouter()
	apiRouter := router.SetupRouter()
	r.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Enable CORS for development
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

// Helper function to remove file
func removeFile(filePath string) error {
	//fmt.Println("di eksekusi")
	if filePath == "" {
		return nil
	}
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
