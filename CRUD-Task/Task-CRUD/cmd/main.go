package main

import (
	"Task-CRUD/config"
	"Task-CRUD/delivery"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL database
db, err := config.InitPostgres(cfg)
if err != nil {
	log.Fatalf("Failed to initialize PostgreSQL: %v", err)
}


	// Initialize Redis
	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize router
	router := delivery.NewRouter(db, redisClient)

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a separate goroutine
	go func() {
		log.Printf("🚀 Server is running on port %s...", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	<-shutdownChan // Menunggu sinyal shutdown
	log.Println("🛑 Shutting down server...")

	// Berikan waktu untuk menyelesaikan koneksi yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("✅ Server successfully stopped")
}
