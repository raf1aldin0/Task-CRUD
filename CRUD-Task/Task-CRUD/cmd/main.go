package main

import (
	"Task-CRUD/config"
	"Task-CRUD/delivery"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL connection
	db, err := config.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal menginisialisasi PostgreSQL: %v", err)
	}
	defer func() {
		if err := config.ClosePostgres(); err != nil {
			log.Printf("⚠️ Gagal menutup koneksi database: %v", err)
		} else {
			log.Println("✅ Koneksi database ditutup dengan aman")
		}
	}()

	// Initialize router
	router := delivery.NewRouter(db)

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
		log.Printf("🚀 Server berjalan di port %s...", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan // Menunggu sinyal shutdown
	log.Println("🛑 Menutup server...")

	// Konteks dengan timeout 10 detik sebelum shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Gagal menghentikan server: %v", err)
	}

	log.Println("✅ Server berhasil dimatikan dengan aman")
}
