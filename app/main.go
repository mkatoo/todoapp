package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkatoo/todoapp/app/handlers"
	"github.com/mkatoo/todoapp/app/middlewares/accesslog"
	"github.com/mkatoo/todoapp/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=postgres password=password dbname=todoapp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Token{}, &models.Task{})
	if err != nil {
		panic("failed to migrate database")
	}

	router := gin.Default()
	router.Use(accesslog.AccessLogMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	handlers.RegisterUserHandler(router, db)
	handlers.RegisterAuthHandler(router, db)
	handlers.RegisterTaskHandler(router, db)

	err = router.Run(":8080")
	if err != nil {
		panic("failed to start server")
	}
}
