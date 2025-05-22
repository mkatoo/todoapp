package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkatoo/todoapp/app/api"
	"github.com/mkatoo/todoapp/app/models"
	"gorm.io/gorm"
)

func RegisterUserHandler(router *gin.Engine, db *gorm.DB) {
	router.GET("/users", func(c *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch users"})
			return
		}
		var usersResponse []api.User
		for _, user := range users {
			usersResponse = append(usersResponse, api.User{
				Id:    int(user.ID),
				Name:  user.Name,
				Email: user.Email,
			})
		}
		c.JSON(http.StatusOK, usersResponse)
	})
	router.POST("/users", func(c *gin.Context) {
		var request api.UserCreateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, api.Error{Error: "invalid input"})
			return
		}
		exists, err := models.IsExists(db, request.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to check user existence"})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, api.Error{Error: "user already exists"})
			return
		}
		user, err := models.NewUser(request.Name, request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to create user"})
			return
		}
		if err := db.Create(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to create user"})
			return
		}
		userResponse := api.User{
			Id:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
		c.JSON(http.StatusCreated, userResponse)
	})
}
