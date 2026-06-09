package handlers


import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	fmt.Println("Home", c.RemoteIP())
	// fmt.Println(r.RemoteAddr)
	c.JSON(200, gin.H{
		"message": "welcome" + c.ClientIP(),
	})
}



func NotFound(c *gin.Context) {
	if c.Request == nil {
		fmt.Println("Request is nil")
		return
	}
	// c.HTML(
	// 	http.StatusNotFound,
	// 	"404.html",
	// 	gin.H{
	// 		"title": "Page not found",
	// 	},
	// )
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})


	// w.Header().Set("Content-Type", "text/html")
	// w.WriteHeader(http.StatusNotFound)
	// fmt.Fprintln(w, `
	// 	<h1>404</h1>
	// 	<p>Oops! Page not found</p>
	// `)

	// fmt.Println("404")
	// fmt.Println(w,r)
}
