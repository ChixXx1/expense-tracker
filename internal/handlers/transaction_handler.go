package handlers

import (
	"net/http"
	"time"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	storage database.Storage
}

func(h *TransactionHandler) GetTransaction(ctx *gin.Context) {
	filters := database.TransactionFilters{}

	if startDateStr := ctx.Query("start_date"); startDateStr != ""{
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start_date format, use YYYY-MM-DD",
			})
			return
		}
		filters.StartDate = &startDate
	} 

	if endDateStr := ctx.Query("end_date"); endDateStr != ""{
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid end_date format, use YYYY-MM-DD",
			})
			return
		}
		filters.EndDate = &endDate
	}



	transaction, err := h.storage.GetTransactions(filters)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get transactions",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
	})
}