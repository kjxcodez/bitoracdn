package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"bitoracdn/origin/db"
)

func main() {
	// Load .env automatically
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found, using system env vars")
    }
	// Connect DB
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer dbConn.Close()

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := dbConn.Pool.Ping(ctx); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("db disconnected"))
			return
		}
		w.Write([]byte("ok"))
	})

	fmt.Println("✅ Origin server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
