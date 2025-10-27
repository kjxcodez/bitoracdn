package main

import (
	"log"
	"net/http"
	"os"

	"bitoracdn/origin/db"
	"bitoracdn/origin/handlers"
	"bitoracdn/origin/storage"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found")
	}

	// Connect DB
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer dbConn.Close()

	// Run migrations
	// if err := db.RunMigrations(); err != nil {
	// 	log.Fatalf("Migration error: %v", err)
	// }


	// Initialize Supabase S3 client
	storageClient := storage.NewSupabaseStorage(
		os.Getenv("SUPABASE_ENDPOINT"),
		os.Getenv("SUPABASE_BUCKET"),
		os.Getenv("SUPABASE_ACCESS_KEY"),
	)

	if err != nil {
		log.Fatalf("Storage init error: %v", err)
	}

	// Setup handlers
	uploader := &handlers.UploadHandler{
		DB:      dbConn,
		Storage: storageClient,
		Bucket:  os.Getenv("SUPABASE_BUCKET"),
	}

	r := mux.NewRouter()
	r.HandleFunc("/upload/init", uploader.UploadInitHandler).Methods("POST")
	r.HandleFunc("/upload/complete", uploader.UploadCompleteHandler).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("✅ Origin server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
