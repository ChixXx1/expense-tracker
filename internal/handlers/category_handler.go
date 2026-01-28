package handlers

import (
	"net/http"
	"strconv"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/ChixXx1/expense-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	storage database.Storage
}

func NewCategoryHandler(storage database.Storage) *CategoryHandler {
	return &CategoryHandler{
		storage: storage,
	}
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	categories, err := h.storage.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get categories",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func (h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category ID",
		})
		return
	}

	category, err := h.storage.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "failed to get category by ID",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"category": category,
	})

}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var category models.Category

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := category.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.storage.CreateCategory(&category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create category: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "category created successfully",
		"category": category,
	})
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category ID",
		})
		return
	}

	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	category.ID = id

	if err := category.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.storage.UpdateCategory(&category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update category: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "category updated successfully",
		"category": category,
	})
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category ID",
		})
		return
	}

	if err := h.storage.DeleteCategory(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete category: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}
