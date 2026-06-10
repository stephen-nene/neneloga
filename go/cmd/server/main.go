package main

import (
	// "net/http"
	"os"
	"fmt"
	"log"
	"context"

	"neneloga/internal/router"
	"neneloga/internal/db"
)

func main() {
	fmt.Print("Starting server...")

	// r := gin.Default()

	r := router.SetupRouter()

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	r.SetTrustedProxies([]string{"192.168.1.2"})

	proxy := os.Getenv("TRUSTED_PROXY") // "192.168.1.2" in dev, real LB IP in prod

	if proxy != "" {
		r.SetTrustedProxies([]string{proxy})
	}

	// r.POST("/rait", server.FileWrite)


	// log.Fatal(
	// http.ListenAndServe(":8080", nil)
	fmt.Println("server running on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		fmt.Println(err)
	}
	// )

}
