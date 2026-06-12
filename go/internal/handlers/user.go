package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	
	"neneloga/internal/db"
	"neneloga/internal/models"
)

// GetUsers godoc
// @Summary      Get all users
// @Description  Get a list of all users and their logs
// @Tags         users
// @Produce      json
// @Success      200  {array}   models.User
// @Router       /users/ [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Preload("Logs").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get a user
// @Description  Get a single user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := db.DB.Preload("Logs").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User data"
// @Success      201   {object}  models.User
// @Router       /users/ [post]
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "User ID"
// @Param        user  body      models.User  true  "User data"
// @Success      200   {object}  models.User
// @Router       /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields (GORM's Updates method only updates non-zero fields by default)
	db.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	db.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
