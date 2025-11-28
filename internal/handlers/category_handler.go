package handlers

import (
	"net/http"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	storage database.Storage
}

func NewCategoryHandler(storage database.Storage) *CategoryHandler{
	return &CategoryHandler{
		storage: storage,
	}
}

func(h *CategoryHandler) GetCategories(ctx *gin.Context){
	categories, err := h.storage.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get categories",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}