package handlers

// import (
// 	"net/http"
// 	"github.com/gin-gonic/gin"

// )

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"neneloga/internal/db"
	"neneloga/internal/models"
)

// ServerHandler struct to hold dependencies
type ServerHandler struct {
	DB *gorm.DB
}

// NewServerHandler creates a new server handler
func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		DB: db.DB,
	}
}

// Pagination helper
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// GetServers godoc
// @Summary      Get all servers
// @Description  Get a list of servers with filtering, sorting, and pagination
// @Tags         servers
// @Produce      json
// @Param        status     query     string  false  "Filter by status: active, deleted, offline"
// @Param        os         query     string  false  "Filter by OS"
// @Param        user_id    query     int     false  "Filter by user ID"
// @Param        search     query     string  false  "Search in name, hostname, IP"
// @Param        sort       query     string  false  "Sort by field (prefix with - for DESC)"  Enums(-created_at, created_at, -name, name, -status, status)
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        page_size  query     int     false  "Items per page (default: 10, max: 100)"
// @Success      200        {object}  map[string]interface{}
// @Router       /servers [get]
func (h *ServerHandler) GetServers(c *gin.Context) {
	var servers []models.Server

	// Start building our query
	query := h.DB.Model(&models.Server{})

	// === FILTERING ===

	// Filter by status - this is our enum field
	if status := c.Query("status"); status != "" {
		// Only apply status filter if it's a valid status
		if models.ServerStatus(status).IsValid() {
			query = query.Where("status = ?", status)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid status value",
				"valid_values": models.AllStatuses(),
			})
			return
		}
	} else {
		// By default, don't show deleted items in listing
		// Remove this if you want to show all by default
		query = query.Where("status != ?", models.StatusDeleted)
	}

	// Filter by operating system
	if os := c.Query("os"); os != "" {
		query = query.Where("os = ?", os)
	}

	// Filter by user ID
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uint(id))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
	}

	// Search across multiple fields (name, hostname, IP address)
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%" // MySQL/PostgreSQL LIKE pattern
		query = query.Where(
			"name LIKE ? OR hostname LIKE ? OR ip_address LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// === SORTING ===

	// Get sort parameter, default to newest first
	sort := c.DefaultQuery("sort", "-created_at")

	// Parse the sort parameter
	sortField := sort
	sortOrder := "ASC"

	// Check if we're sorting descending (starts with -)
	if strings.HasPrefix(sort, "-") {
		sortField = strings.TrimPrefix(sort, "-")
		sortOrder = "DESC"
	}

	// Whitelist of allowed sort fields (prevents SQL injection)
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
		"status":     true,
		"os":         true,
		"id":         true,
	}

	if !allowedSortFields[sortField] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sort field. Allowed: " + strings.Join(getAllowedSortFields(allowedSortFields), ", "),
		})
		return
	}

	query = query.Order(sortField + " " + sortOrder)

	// === PAGINATION ===

	// Parse page number (default: 1)
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse page size (default: 10, max: 100)
	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			if ps > 100 {
				ps = 100 // Cap at 100
			}
			pageSize = ps
		}
	}

	// Count total records (for pagination info)
	var total int64
	// We need to count before pagination
	countQuery := *query // Copy the query
	countQuery.Count(&total)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// === EXECUTE QUERY ===

	// Preload relationships if needed
	// query = query.Preload("Logs")
	// query = query.Preload("User")

	// Execute the query
	if err := query.Find(&servers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch servers"})
		return
	}

	// === RETURN RESPONSE ===

	// Calculate total pages
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"data": servers,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
			"has_more":    page < totalPages,
		},
	})
}

// GetServer godoc
// @Summary      Get a single server
// @Description  Get server by ID with its logs
// @Tags         servers
// @Produce      json
// @Param        id   path      int  true  "Server ID"
// @Success      200  {object}  models.Server
// @Router       /servers/{id} [get]
func (h *ServerHandler) GetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var server models.Server

	// Preload logs for the detail view
	if err := h.DB.Preload("Logs").First(&server, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server"})
		return
	}

	c.JSON(http.StatusOK, server)
}

// CreateServer godoc
// @Summary      Create a new server
// @Description  Add a new server to the database
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        server  body      models.Server  true  "Server object"
// @Success      201     {object}  models.Server
// @Router       /servers [post]
func (h *ServerHandler) CreateServer(c *gin.Context) {
	var server models.Server

	// Bind JSON body to server struct
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default status if not provided
	if server.Status == "" {
		server.Status = models.StatusActive
	}

	// Validate status
	if !server.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid status",
			"valid_values": models.AllStatuses(),
		})
		return
	}

	// Create the server
	if err := h.DB.Create(&server).Error; err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Server with this name, IP, or hostname already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create server"})
		return
	}

	c.JSON(http.StatusCreated, server)
}

// UpdateServer godoc
// @Summary      Update a server
// @Description  Update an existing server
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        id      path      int            true  "Server ID"
// @Param        server  body      models.Server  true  "Server object"
// @Success      200     {object}  models.Server
// @Router       /servers/{id} [put]
func (h *ServerHandler) UpdateServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var server models.Server

	// Check if server exists
	if err := h.DB.First(&server, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server"})
		return
	}

	// Bind JSON body
	var update models.Server
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status if provided
	if update.Status != "" && !update.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid status",
			"valid_values": models.AllStatuses(),
		})
		return
	}

	// Update the server (don't update zero values unless explicitly set)
	// Using map for partial updates
	updates := map[string]interface{}{}

	if update.Name != "" {
		updates["name"] = update.Name
	}
	if update.IPAddress != "" {
		updates["ip_address"] = update.IPAddress
	}
	if update.Hostname != "" {
		updates["hostname"] = update.Hostname
	}
	if update.Os != "" {
		updates["os"] = update.Os
	}
	if update.Status != "" {
		updates["status"] = update.Status
	}
	if update.UserID != 0 {
		updates["user_id"] = update.UserID
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.DB.Model(&server).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update server"})
		return
	}

	c.JSON(http.StatusOK, server)
}

// DeleteServer godoc
// @Summary      Soft delete a server
// @Description  Change server status to 'deleted' instead of hard delete
// @Tags         servers
// @Produce      json
// @Param        id   path      int  true  "Server ID"
// @Success      200  {object}  gin.H
// @Router       /servers/{id} [delete]
func (h *ServerHandler) DeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var server models.Server

	// Check if server exists
	if err := h.DB.First(&server, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server"})
		return
	}

	// Check if already deleted
	if server.Status == models.StatusDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server is already deleted"})
		return
	}

	// Soft delete - just change the status
	if err := h.DB.Model(&server).Update("status", models.StatusDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Server deleted successfully",
		"server_id": server.ID,
	})
}

// Helper function to get allowed sort fields as slice
func getAllowedSortFields(fields map[string]bool) []string {
	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}
	return result
}
