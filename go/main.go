package main

import (
	"net/http"

	"neneloga.local/server"

	// "log"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("Starting server...")

	r := gin.Default()
  r.SetTrustedProxies([]string{"192.168.1.2"})

	r.GET("/ping",server.Ping )
	




	http.HandleFunc("/",server.Home)
	http.HandleFunc("/rait",server.FileWrite)

	// 404
	http.HandleFunc("/404",server.NotFound)



	// log.Fatal(
		// http.ListenAndServe(":8080", nil)
		r.Run()
		fmt.Println("server running on :8080")
	// )

}
