package main

import (
	"log"
	"net/http"
	"os"
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
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

	// Get all devices from the database
	devices, err := storeContainer.GetAllDevices()
	if err != nil {
		log.Fatalf("Failed to get devices from database: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		client.AddEventHandler(controllers.EventHandler)

		if client.Store.ID == nil {
			// New login
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				log.Fatalf("Failed to connect client: %v", err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					// Render QR code
					fmt.Println("QR code:", evt.Code)
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		} else {
			// Already logged in, just connect
			err = client.Connect()
			if err != nil {
				log.Fatalf("Failed to connect client: %v", err)
			}
		}
	}

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
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, corsHandler))
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		client.Disconnect()
	}
}
