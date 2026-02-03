package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/ChixXx1/expense-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	storage database.Storage
}

func NewBudgetHandler(storage database.Storage) *BudgetHandler {
	return &BudgetHandler{
		storage: storage,
	}
}

func (h *BudgetHandler) GetBudgets(ctx *gin.Context) {
	filters := database.BudgetFilters{}

	if categoryIDStr := ctx.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil || categoryID <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "category_id must be a positive integer",
			})
			return
		}
		filters.CategoryID = &categoryID
	}

	if period := ctx.Query("period"); period != "" {
		validPeriods := map[string]bool{
			models.BudgetPeriodMonthly: true,
			models.BudgetPeriodWeekly:  true,
			models.BudgetPeriodYearly:  true,
		}
		if !validPeriods[period] {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "period must be 'monthly', 'weekly' or 'yearly'",
			})
			return
		}
		filters.Period = &period
	}

	if monthStr := ctx.Query("month"); monthStr != "" {
		month, err := time.Parse("2006-01", monthStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid month format, use YYYY-MM",
			})
			return
		}
		filters.Month = &month
	}

	budgets, err := h.storage.GetBudgets(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get budgets",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"budgets": budgets,
		"count":   len(budgets),
	})
}

func (h *BudgetHandler) GetBudgetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid budget ID",
		})
		return
	}

	budget, err := h.storage.GetBudgetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "budget not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"budget": budget,
	})
}

func (h *BudgetHandler) CreateBudget(ctx *gin.Context) {
	var budget models.Budget

	if err := ctx.ShouldBindJSON(&budget); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := budget.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.storage.CreateBudget(&budget); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create budget: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "budget created successfully",
		"budget":  budget,
	})
}

func (h *BudgetHandler) UpdateBudget(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid budget ID",
		})
		return
	}

	var budget models.Budget
	if err := ctx.ShouldBindJSON(&budget); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	budget.ID = id

	if err := budget.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.storage.UpdateBudget(&budget); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update budget: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "budget updated successfully",
		"budget":  budget,
	})
}

func (h *BudgetHandler) DeleteBudget(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid budget ID",
		})
		return
	}

	if err := h.storage.DeleteBudget(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete budget: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "budget deleted successfully",
	})
}
