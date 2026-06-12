package main

import (
    "fmt"
    "log"
    "os"

    "neneloga/internal/db"
    "neneloga/internal/redis"
    "neneloga/internal/router"
)

func main() {
    fmt.Println("Starting server...")

    // Setup router
    r := router.SetupRouter()

    // Connect to database
    if err := db.Connect(); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Get the underlying sql.DB for connection pooling
    sqlDB, err := db.DB.DB()
    if err == nil {
        defer sqlDB.Close()
    }

    // Connect to Redis (optional - log warning instead of fatal if not available)
    if err := redis.Connect(); err != nil {
        log.Printf("Warning: Redis not available: %v", err)
        // Don't fatal - allow server to run without Redis
    }

    // Configure trusted proxies
    proxy := os.Getenv("TRUSTED_PROXY")
    if proxy != "" {
        r.SetTrustedProxies([]string{proxy})
    } else {
        // Default fallback for development
        r.SetTrustedProxies([]string{"192.168.1.2"})
    }

    fmt.Println("Server running on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
