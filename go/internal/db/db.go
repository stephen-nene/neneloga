package db

import (
    // "database/sql"
    "os"
    "fmt"
    "context"
	"log"


	"github.com/jackc/pgx/v5"
"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"neneloga/internal/models"
    // _ "github.com/lib/pq"
)

func Connect2() (*pgx.Conn, error) {
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	// 	os.Getenv("POSTGRES_USER"),
	// 	os.Getenv("POSTGRES_PASSWORD"),
	// 	os.Getenv("POSTGRES_HOST"),
	// 	os.Getenv("POSTGRES_PORT"),
	// 	os.Getenv("POSTGRES_DB"))

    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("failed to load .env file: %w", err)
    }

	// fmt.Println("DB_URL hia:", os.Getenv("DB_URL"))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	   // Verify the connection is actually working
    if err := conn.Ping(context.Background()); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
	// conn.SetMaxOpenConns(25)
	// conn.SetMaxIdleConns(5)
	// conn.SetConnMaxLifetime(time.Hour)
	return conn, nil
}


var DB *gorm.DB

func Connect() error {
	// Attempt to load .env, but don't fail if it doesn't exist (e.g. in production)
	_ = godotenv.Load()

    dsn := os.Getenv("DB_URL")

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
		return err
    }

    DB = database

    return DB.AutoMigrate(
		&models.User{},
        &models.Log{},
    )
}
// func InitDB() *sql.DB {
//     db, _ := sql.Open("postgres", os.Getenv("DB_URL"))
//     return db
// }

// var DB *gorm.DB

// func Connect() (*gorm.DB, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load .env file: %w", err)
// 	}
// 	dsn := os.Getenv("DB_URL")
// 	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
// 	}
// 	DB = database
// 	// This automatically creates/updates the tables based on your structs!
// 	err = DB.AutoMigrate(&models.User{}, &models.Log{})
// 	// DB.Create(&models.User{Username: "Steve", Email: "me@stevenene.top", Password: "password"})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to migrate database: %w", err)
// 	}
// 	return DB, nil
// }
