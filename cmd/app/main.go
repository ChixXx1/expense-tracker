package main

import (
	"net/http"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/ChixXx1/expense-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	storage := database.NewMemoryStorage()
	categoryHandler := handlers.NewCategoryHandler(storage)

	r := gin.Default()


	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "WEB-APPLICATION GO+REACT",
		})
	})
	r.GET("/categories", categoryHandler.GetCategories)
	r.GET("/categories/:id", categoryHandler.GetCategoryByID)

	
	r.Run(":8080")
}