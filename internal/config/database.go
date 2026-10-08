package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "community_waste_db"),
		getEnv("DB_PORT", "5432"),
	)
}

func InitDB() *gorm.DB {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")

	db, err := gorm.Open(postgres.Open(DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("gagal koneksi ke database: %v", err)
	}

	if err := runSQLFile(db, "migrations/0001_init.sql"); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}

	if os.Getenv("RUN_SEED") == "true" {
		if err := runSQLFile(db, "migrations/0002_seed.sql"); err != nil {
			log.Printf("gagal seeding (non-fatal): %v", err)
		} else {
			log.Println("seeding selesai")
		}
	}

	return db
}

func runSQLFile(db *gorm.DB, path string) error {
	abs := path
	if _, err := os.Stat(abs); err != nil {
		alt := filepath.Join("..", "..", path)
		if _, err2 := os.Stat(alt); err2 == nil {
			abs = alt
		} else {
			return err
		}
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return db.Exec(string(b)).Error
}
