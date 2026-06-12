package redis

import (
    "context"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
    Client = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"), // e.g., "localhost:6379"
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := Client.Ping(ctx).Result()
    return err
}
