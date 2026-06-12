package handlers


import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func Home(c *gin.Context) {
// 	fmt.Println("Home", c.RemoteIP())
// 	fmt.Println(c.Request,c.Request.Response)
// 	c.JSON(200, gin.H{
// 		// "message": "welcome" + c.ClientIP(),
// 		// "header": c.Request.Header,
// 		"home": c.Request.RequestURI,
// 		"health": c.Request.RequestURI,
// 		"chuck": c.Request.RequestURI,
// 		"users": c.Request.RequestURI,
// 		"users/:id": c.Request.RequestURI,
// 	})
// }

// Home godoc
// @Summary      Home endpoint
// @Description  Get a list of all available endpoints
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       / [get]
func Home(c *gin.Context) {
	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	baseURL := scheme + "://" + host

	c.JSON(200, gin.H{
		"base_url": baseURL,
		"endpoints": []string{
			baseURL + "/",
			baseURL + "/health",
			baseURL + "/chuck",
			baseURL + "/users",
			baseURL + "/users/:id",
		},
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
