package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkatoo/todoapp/app/api"
	"github.com/mkatoo/todoapp/app/models"
	"gorm.io/gorm"
)

func RegisterTaskHandler(router *gin.Engine, db *gorm.DB) {
	router.GET("/tasks", func(c *gin.Context) {
		user, err := GetUserByToken(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch user"})
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, api.Error{Error: "unauthorized"})
			return
		}
		var tasks []models.Task
		if err := db.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch tasks"})
			return
		}
		tasksResponse := make([]api.Task, 0, len(tasks))
		for _, task := range tasks {
			tasksResponse = append(tasksResponse, api.Task{
				Id:        int(task.ID),
				Completed: task.Completed,
				Content:   task.Content,
				CreatedAt: task.CreatedAt,
				UpdatedAt: task.UpdatedAt,
			})
		}
		c.JSON(http.StatusOK, tasksResponse)
	})
	router.POST("/tasks", func(c *gin.Context) {
		user, err := GetUserByToken(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch user"})
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, api.Error{Error: "unauthorized"})
			return
		}
		var request api.TaskCreateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, api.Error{Error: "invalid input"})
			return
		}
		task := models.Task{
			UserID:    user.ID,
			Content:   request.Content,
			Completed: false,
		}
		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to create task"})
			return
		}
		taskResponse := api.Task{
			Id:        int(task.ID),
			Content:   task.Content,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}
		c.JSON(http.StatusCreated, taskResponse)
	})
	router.PUT("/tasks/:id", func(c *gin.Context) {
		user, err := GetUserByToken(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch user"})
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, api.Error{Error: "unauthorized"})
			return
		}
		var request api.TaskUpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, api.Error{Error: "invalid input"})
			return
		}
		taskID := c.Param("id")
		var task models.Task
		if err := db.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, api.Error{Error: "task not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to fetch task"})
			return
		}
		task.Content = request.Content
		task.Completed = request.Completed
		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, api.Error{Error: "failed to update task"})
			return
		}
		taskResponse := api.Task{
			Id:        int(task.ID),
			Content:   task.Content,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}
		c.JSON(http.StatusOK, taskResponse)
	})
}

func GetUserByToken(c *gin.Context, db *gorm.DB) (*models.User, error) {
	authString := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authString, "Bearer ")
	if tokenString == "" {
		return nil, nil
	}
	var token models.Token
	if err := db.Where("token = ?", tokenString).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var user models.User
	if err := db.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
