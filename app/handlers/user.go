package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkatoo/todoapp/app/models"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterUserHandler(router *gin.Engine, db *gorm.DB) {
	router.GET("/users", func(c *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	router.POST("/users", func(c *gin.Context) {
		var request CreateUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user, err := models.NewUser(request.Name, request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		if err := db.Create(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, user)
	})
}
