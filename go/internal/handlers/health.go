package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// types
type ChuckResponse struct {
	IconURL string `json:"icon_url"`
	Value   string `json:"value"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Server is healthy & aLIVE!"})
}

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
