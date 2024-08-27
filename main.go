package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"wagobot.com/controllers"
	"wagobot.com/db"
	"wagobot.com/middleware"
	"wagobot.com/router"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Memuat nilai dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Gagal memuat file .env: %v", err)
	}

	// Mengambil nilai dari variabel lingkungan
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")
	botName := os.Getenv("BOT_NAME")

	if botName == "" {
		botName = "WINDOWS"
	}

	store.DeviceProps.PlatformType = waProto.DeviceProps_SAFARI.Enum()
	store.DeviceProps.Os = proto.String(botName)	

	// Mengatur logging untuk database
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Menginisialisasi SQL store
	storeContainer, err := sqlstore.New("sqlite3", dbPath, dbLog)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}	

	// Menetapkan storeContainer ke variabel package controller
	controllers.SetStoreContainer(storeContainer)

	// Mengambil semua perangkat dari database
	devices, err := storeContainer.GetAllDevices()
	if err != nil {
		log.Fatalf("Gagal mengambil perangkat dari database: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		controllers.AddClient(client)
	}

	// Inisialisasi database
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	defer db.CloseDB()

	// Mengatur router
	r := router.SetupRouter()
	apiRouter := router.SetupRouter()
	r.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))

	//Tambahkan Endpoin swager
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Menyajikan file statis
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Mengaktifkan CORS untuk pengembangan
	corsHandler := middleware.SetupCORS(r)

	log.Printf("Server berjalan di port %s\n", port)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, corsHandler))
	}()

	// Shutdown yang aman
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	for _, device := range devices {
		client := whatsmeow.NewClient(device, clientLog)
		client.Disconnect()
	}
}
