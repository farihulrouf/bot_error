package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv" // Import godotenv package
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"wagobot.com/controllers"
	"wagobot.com/db"
	"wagobot.com/router"
)

var client *whatsmeow.Client

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Initialize database connection
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}
	if err := db.InitDB(dbPath); err != nil {
		panic(err)
	}
	defer db.CloseDB()

	// Initialize WhatsApp connection
	container, err := sqlstore.New("sqlite3", dbPath, dbLog)
	if err != nil {
		panic(err)
	}
	defer container.Close()

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(controllers.EventHandler)
	controllers.SetClient(client)
	controllers.ScanQrCode(client)
	//controllers.EventHandler() // Call EventHandler here
	controllers.EventHandler(client)

	// Setup router with client
	r := router.SetupRouter(client)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	go func() {
		address := fmt.Sprintf(":%s", port)
		log.Printf("Server running on port %s", port)
		if err := http.ListenAndServe(address, r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
