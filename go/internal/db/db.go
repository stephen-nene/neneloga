package db

import (
    "os"
    "fmt"
    "context"
    "log"
    "strings"

    "github.com/jackc/pgx/v5"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "neneloga/internal/models"
)

func Connect2() (*pgx.Conn, error) {
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
// type DB struct {
// 	*gorm.DB
// }

func Connect() error {
	// Attempt to load .env, but don't fail if it doesn't exist (e.g. in production)
	_ = godotenv.Load()

    dsn := os.Getenv("DB_URL")

    var dialector gorm.Dialector
    if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "host=") {
        dialector = postgres.Open(dsn)
    } else {
        // Fallback to SQLite
        dialector = sqlite.Open("neneloga.db")
    }

    database, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
		return err
    }

    DB = database
	err = DB.AutoMigrate(&models.User{}, &models.Log{}, &models.Server{})
	if err != nil {
		return err
	}

	// 👇 seed only if empty
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		DB.Create(&models.User{
			Username: "Steve",
			Email:    os.Getenv("ADMIN_EMAIL"),
			Password: os.Getenv("ADMIN_PASSWORD"),
		})
	}

	return nil
    // return DB.AutoMigrate(
	// 	&models.User{},
    //     &models.Server{},
    //     &models.Log{},
    // )
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
