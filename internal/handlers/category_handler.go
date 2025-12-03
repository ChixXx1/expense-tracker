package handlers

import (
	"net/http"
	"strconv"

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

func(h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	idParam := ctx.Param(":id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "inavalid category ID",
		})
		return
	}

	category, err := h.storage.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "failed to get category by ID",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

func(h *CategoryHandler) CreateCategory(ctx *gin.Context) error{



	return nil
}