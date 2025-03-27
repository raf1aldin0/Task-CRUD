package config

import (
	"database/sql" // Gunakan sql.DB untuk connection pooling
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Variabel global untuk database pooling
var (
	db    *gorm.DB
	sqlDB *sql.DB
	once  sync.Once // Untuk memastikan koneksi dibuat hanya sekali
)

// InitPostgres menginisialisasi koneksi database dengan pooling
func InitPostgres(cfg *Config) (*gorm.DB, error) {
	var err error

	// Gunakan sync.Once untuk memastikan koneksi hanya dibuat sekali
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("❌ Gagal menghubungkan ke PostgreSQL: %v", err)
			return
		}

		// Ambil sql.DB untuk connection pooling
		sqlDB, err = db.DB()
		if err != nil {
			log.Printf("❌ Gagal mendapatkan sql.DB dari GORM: %v", err)
			return
		}

		// Konfigurasi pooling
		sqlDB.SetMaxOpenConns(25)        // Maksimum koneksi terbuka
		sqlDB.SetMaxIdleConns(10)        // Maksimum koneksi idle
		sqlDB.SetConnMaxLifetime(5 * 60) // Waktu maksimum koneksi hidup dalam detik

		log.Println("✅ Koneksi PostgreSQL berhasil dibuat dengan pooling")
	})

	return db, err
}

// ClosePostgres menutup koneksi database dengan aman
func ClosePostgres() error {
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}
