package db

import (
    // "database/sql"
    "os"
    "fmt"
    "context"
	"log"


	"github.com/jackc/pgx/v5"
"github.com/joho/godotenv"
    // _ "github.com/lib/pq"
)

func Connect() (*pgx.Conn, error) {
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


// func InitDB() *sql.DB {
//     db, _ := sql.Open("postgres", os.Getenv("DB_URL"))
//     return db
// }
