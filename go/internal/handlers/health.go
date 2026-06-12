package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "sync"
    "time"

    "neneloga/internal/db"
    "neneloga/internal/redis"

    "github.com/gin-gonic/gin"
)

type ChuckResponse struct {
	IconURL string `json:"icon_url"`
	Value   string `json:"value"`
}

// type ServiceStatus struct {
//     Name    string `json:"name"`
//     Status  string `json:"status"` // "UP" or "DOWN"
//     Url     *string `json:"url,omitempty"`
//     Message string `json:"message,omitempty"`
//     Latency string `json:"latency,omitempty"`
// }
type ServiceStatus struct {
    Name    string `json:"name"`
    Url     string `json:"url,omitempty"`     // Will be omitted if empty
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`  // Will be omitted if empty
    Latency string `json:"latency,omitempty"`  // Will be omitted if empty
}

type HealthResponse struct {
    Status   string          `json:"status"` // "UP" or "DEGRADED" or "DOWN"
    Server   ServiceStatus   `json:"server"`
    Services []ServiceStatus `json:"services,omitempty"`
}

func Health(c *gin.Context) {
    var (
        wg      sync.WaitGroup
        mu      sync.Mutex
        results []ServiceStatus
    )

    // Server status (always UP if we can serve this request)
    serverStatus := ServiceStatus{
        Name:   "server",
        Status: "UP",
        Message: "Server is running",
    }

    // External services to check
    services := []struct {
        Name string
        URL  string
    }{
        {Name: "github", URL: "https://api.github.com"},
        {Name: "safaricom", URL: os.Getenv("SAFARICOM_API_URL")}, // e.g., "https://api.safaricom.co.ke"
        {Name: "buni", URL: os.Getenv("BUNI_API_URL")},           // e.g., "https://api.buni.com"
    }

    // Check Database
    wg.Add(1)
    go func() {
        defer wg.Done()
        dbStatus := checkDatabase()
        mu.Lock()
        results = append(results, dbStatus)
        mu.Unlock()
    }()

    // Check Redis (if configured)
    wg.Add(1)
    go func() {
        defer wg.Done()
        redisStatus := checkRedis()
        mu.Lock()
        results = append(results, redisStatus)
        mu.Unlock()
    }()

    // Check external services
    for _, svc := range services {
        // fmt.Println(svc.URL)
        if svc.URL == "" {
            continue // Skip if URL not configured
        }
        wg.Add(1)
        go func(name, url string) {
            defer wg.Done()
            extStatus := checkExternalService(name, url)
            mu.Lock()
            results = append(results, extStatus)
            mu.Unlock()
        }(svc.Name, svc.URL)
    }

    wg.Wait()

    // Determine overall status
    overallStatus := "UP"
    for _, svc := range results {
        if svc.Status == "DOWN" {
            overallStatus = "DEGRADED"
            break
        }
    }
    if serverStatus.Status == "DOWN" {
        overallStatus = "DOWN"
    }

    c.JSON(http.StatusOK, HealthResponse{
        Status: overallStatus,
        Server: serverStatus,
        Services: results,
    })
}

func checkDatabase() ServiceStatus {
    start := time.Now()

    sqlDB, err := db.DB.DB()
    if err != nil {
        return ServiceStatus{
            Name:    "database",
            Status:  "DOWN",
            Message: fmt.Sprintf("Failed to get database instance: %v", err),
        }
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if err := sqlDB.PingContext(ctx); err != nil {
        return ServiceStatus{
            Name:    "database",
            Status:  "DOWN",
            Message: fmt.Sprintf("Database ping failed: %v", err),
        }
    }

    latency := time.Since(start)
    return ServiceStatus{
        Name:    "database",
        Status:  "UP",
        Message: "Database connection is healthy",
        Latency: latency.String(),
    }
}

func checkRedis() ServiceStatus {
    if redis.Client == nil {
        return ServiceStatus{
            Name:    "redis",
            Status:  "NOT_CONFIGURED",
            Message: "Redis client not initialized",
        }
    }

    start := time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    _, err := redis.Client.Ping(ctx).Result()
    if err != nil {
        return ServiceStatus{
            Name:    "redis",
            Status:  "DOWN",
            Message: fmt.Sprintf("Redis ping failed: %v", err),
        }
    }

    latency := time.Since(start)
    return ServiceStatus{
        Name:    "redis",
        Status:  "UP",
        Message: "Redis connection is healthy",
        Latency: latency.String(),
    }
}

func checkExternalService(name, url string) ServiceStatus {
    start := time.Now()

    // Create status without URL first
    status := ServiceStatus{
        Name: name,
    }

    // Only include URL if it's not empty
    // if url != "" {
    //     status.Url = url
    // }

    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get(url)
    if err != nil {
        status.Status = "DOWN"
        status.Message = fmt.Sprintf("Failed to reach service: %v", err)
        return status
    }
    defer resp.Body.Close()

    status.Latency = time.Since(start).String()

    status.Status = "UP"
    status.Message = fmt.Sprintf("Service responded with HTTP %d", resp.StatusCode)

    // resp, err := client.Get(url)
    // fmt.Println(resp)
    // if err != nil {
    //     status.Status = "DOWN"
    //     status.Message = fmt.Sprintf("Failed to reach service: %v", err)
    //     return status
    // }
    // defer resp.Body.Close()

    // if resp.StatusCode >= 200 && resp.StatusCode < 300 {
    //     status.Status = "UP"
    //     status.Message = fmt.Sprintf("Service is reachable (HTTP %d)", resp.StatusCode)
    // } else {
    //     status.Status = "DOWN"
    //     status.Message = fmt.Sprintf("Service returned unexpected status: HTTP %d", resp.StatusCode)
    // }

    return status
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
