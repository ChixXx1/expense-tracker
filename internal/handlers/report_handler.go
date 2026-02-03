package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	storage database.Storage
}

func NewReportHandler(storage database.Storage) *ReportHandler {
	return &ReportHandler{
		storage: storage,
	}
}

func (h *ReportHandler) GetFinancialSummary(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date format, use YYYY-MM-DD",
		})
		return
	}

	if startDate.After(endDate) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date must be before end_date",
		})
		return
	}

	summary, err := h.storage.GetFinancialSummary(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get financial summary",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}

func (h *ReportHandler) GetCategorySummary(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date format, use YYYY-MM-DD",
		})
		return
	}

	if startDate.After(endDate) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date must be before end_date",
		})
		return
	}

	summaries, err := h.storage.GetCategorySummary(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get category summary",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"category_summary": summaries,
		"count":            len(summaries),
	})
}

func (h *ReportHandler) GetBudgetReport(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid budget ID",
		})
		return
	}

	report, err := h.storage.GetBudgetReport(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "budget report not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"budget_report": report,
	})
}
