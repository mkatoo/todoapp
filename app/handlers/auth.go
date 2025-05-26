package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkatoo/todoapp/app/api"
	"github.com/mkatoo/todoapp/app/models"
	"gorm.io/gorm"
)

func RegisterAuthHandler(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth", func(c *gin.Context) {
		var request api.AuthRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, api.Error{Error: "invalid input"})
			return
		}
		exists, err := models.IsUserExists(db, request.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to check user existence"})
			return
		}
		if !exists {
			c.JSON(http.StatusUnauthorized, api.Error{Error: "invalid credentials"})
			return
		}

		var user *models.User
		if err = db.Where("email = ?", request.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch user"})
			return
		}
		if !user.CheckPassword(request.Password) {
			c.JSON(http.StatusUnauthorized, api.Error{Error: "invalid credentials"})
			return
		}

		token, err := models.FindOrCreateToken(db, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, api.Token{Token: token.Token})
	})
}
