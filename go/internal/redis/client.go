package redis

import (
    "context"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"), // e.g., "localhost:6379"
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := client.Ping(ctx).Result()
    if err != nil {
        client.Close()
        return err
    }

    Client = client
    return nil
}
