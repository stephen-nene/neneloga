package server

import (
	"fmt"
	// "os"
	"net/http"
	// "log"
	"github.com/gin-gonic/gin"
	)


type Input struct {
	Name string `json:"name"`
	Body []byte `json:"body"`

}

type Output struct {
	Message string `json:"message"`
}

func Ping(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}

func FileWrite(w http.ResponseWriter,r *http.Request){

}


func Home(w http.ResponseWriter,r *http.Request){
	fmt.Println("Home")
	fmt.Println(r.RemoteAddr)
}

func NotFound(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type","application/json")
	fmt.Fprintln(w,`
		<h1>404</h1>
		<p>Oops! Page not found</p>
	`)
	fmt.Println("404")
	// fmt.Println(w,r)
}


