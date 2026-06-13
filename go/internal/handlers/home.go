package handlers


import (
    "io"
	"fmt"
	"net/http"
    "encoding/json"


	"github.com/gin-gonic/gin"
)


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
			// baseURL + "/",
			baseURL + "/health",
			baseURL + "/chuck",
			baseURL + "/users",
			baseURL + "/users/:id",
			baseURL + "/swagger/index.html",
			// baseURL + "/swagger/*any",
		},
	})
}

// ChuckNorris godoc
// @Summary      Chuck Norris Joke
// @Description  Get a random Chuck Norris joke
// @Tags         fun
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /chuck [get]
func ChuckNorris(c *gin.Context) {
	res, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "ERROR", "message": err.Error()})
		return
	}

	formart := c.Query("fmt")
	// https://api.chucknorris.io/jokes/random
	// fmt.Println(res)
	defer res.Body.Close()

	// Read the raw body from the external API
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERROR",
			"message": "Failed to read response",
		})
		return
	}

	var joke ChuckResponse
	if err := json.Unmarshal(body, &joke); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERROR",
			"message": "Failed to parse JSON",
		})
		return
	}

	if formart != "html" {

		// c.Data(res.StatusCode, "application/json", body)
		// default JSON response (cleaned)
		c.JSON(http.StatusOK, gin.H{
			"icon_url": joke.IconURL,
			"joke":     joke.Value,
		})
		return
	}
	// if format == "html" {
	html := fmt.Sprintf(`
			<div>
				<img src="%s" width="100"/>
				<p>%s</p>
				<button onclick="location.reload()">Reload</button>
			</div>
		`, joke.IconURL, joke.Value)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))

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
