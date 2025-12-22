package main

import (
	"net/http"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/ChixXx1/expense-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	//storage := database.NewMemoryStorage()
	storage := database.NewJSONStorage("./data.json")
	categoryHandler := handlers.NewCategoryHandler(storage)

	r := gin.Default()

	r.StaticFile("/favicon.ico", "./static/favicon.ico")
  /* r.GET("/favicon.ico", func(c *gin.Context) {
    c.Status(http.StatusNoContent)
  }) */

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
	r.POST("/categories", categoryHandler.CreateCategory)
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)
	r.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	r.Run(":8080")
}