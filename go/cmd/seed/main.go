package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"neneloga/internal/db"
	"neneloga/internal/models"
)

func main() {
	// Define command line flags
	numUsers := flag.Int("users", 1, "Number of users to create")
	numServers := flag.Int("servers", 2, "Number of servers per user")
	numLogs := flag.Int("logs", 5, "Number of logs per server")

	flag.Parse()

	// Connect to database
	fmt.Println("Connecting to database...")
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Printf("Starting seed: %d users, %d servers per user, %d logs per server...\n", *numUsers, *numServers, *numLogs)

	for i := 1; i <= *numUsers; i++ {
		// Create a user
		password := fmt.Sprintf("password%d", i)
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := models.User{
			Username: fmt.Sprintf("seeduser%d", i),
			Email:    fmt.Sprintf("seeduser%d@example.com", i),
			Password: string(hashed),
			Role:     "user",
		}

		if err := db.DB.Create(&user).Error; err != nil {
			fmt.Printf("Failed to create user %d: %v\n", i, err)
			continue
		}
		fmt.Printf("Created user: %s (password: %s)\n", user.Username, password)

		for j := 1; j <= *numServers; j++ {
			// Create a server for the user
			server := models.Server{
				UserID:    user.ID,
				Name:      fmt.Sprintf("Server-%d-%d", i, j),
				IPAddress: fmt.Sprintf("192.168.%d.%d", i, j),
				Hostname:  fmt.Sprintf("host-%d-%d.local", i, j),
				Os:        "Ubuntu 22.04",
				Status:    models.StatusActive,
			}

			if err := db.DB.Create(&server).Error; err != nil {
				fmt.Printf("  Failed to create server %d: %v\n", j, err)
				continue
			}
			fmt.Printf("  Created server: %s (IP: %s)\n", server.Name, server.IPAddress)

			for k := 1; k <= *numLogs; k++ {
				// Create logs for the server
				logEntry := models.Log{
					UserID:   &user.ID,
					ServerID: server.ID,
					Level:    "info",
					Message:  fmt.Sprintf("This is log entry %d for server %s", k, server.Name),
				}

				if err := db.DB.Create(&logEntry).Error; err != nil {
					fmt.Printf("    Failed to create log %d: %v\n", k, err)
					continue
				}
			}
			fmt.Printf("    Created %d logs for server %s\n", *numLogs, server.Name)
		}
	}

	fmt.Println("Seeding completed successfully!")
	os.Exit(0)
}
