package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChixXx1/expense-tracker/internal/database"
	"github.com/ChixXx1/expense-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	storage database.Storage
}

func NewTransactionHandler(storage database.Storage) *TransactionHandler{
	return &TransactionHandler{
		storage: storage,
	}
}

func(h *TransactionHandler) GetTransactions(ctx *gin.Context) {
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

	if filters.StartDate != nil && filters.EndDate != nil{
		if filters.StartDate.After(*filters.EndDate){
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "start_date must be before end_date",
			})
			return
		}
	}

	if categoryIDStr := ctx.Query("category_id"); categoryIDStr != ""{
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil || categoryID <= 0{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "category_id must be a positive integer",
			})
			return
		}
		filters.CategoryID = &categoryID
	}

	if txTypeStr := ctx.Query("type"); txTypeStr != ""{
		if txTypeStr != models.TransactionTypeIncome && txTypeStr != models.TransactionTypeExpense{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "type must be `income` or `expense`",
			})
			return
		}
		filters.Type = &txTypeStr
	}

	if paymentMethod := ctx.Query("payment_method"); paymentMethod != ""{
		validMetods := map[string]bool{
			models.PaymentMethodCash: true,
			models.PaymentMethodCard: true,
			models.PaymentMethodTransfer: true,
		}
		if !validMetods[paymentMethod]{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "payment_method must be a `cash`, `card` or `transfer`",
			})
			return
		}
		filters.PaymentMethod = &paymentMethod
	}

	if limitStr := ctx.Query("limit"); limitStr != ""{
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be a positive integer",
			})
			return
		}
		filters.Limit = &limit
	}

	if offsetStr := ctx.Query("offset"); offsetStr != ""{
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "offset must be a positive",
			})
			return
		}
		filters.Offset = &offset
	}

	transactions, err := h.storage.GetTransactions(filters)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get transactions",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count": len(transactions), 
		"filters_applied": gin.H{
			"has_start_date": filters.StartDate != nil,
			"has_end_date": filters.EndDate != nil,
			"has_category": filters.CategoryID != nil,
			"has_type": filters.Type != nil,
		},
	})
}

func(h *TransactionHandler) GetTransactionByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	transaction, err := h.storage.GetTransactionByID(id)
	if err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "failed to get transaction by ID",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
	})
}

func(h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	var transaction models.Transaction

	if err := ctx.ShouldBindJSON(&transaction); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if transaction.Amount <= 0{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "transaction amount must be positive",
		})
		return
	}

	if transaction.Type != "income" && transaction.Type != "expense"{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "transaction type must be 'income' or 'expense'",
		})
		return
	}

	if transaction.PaymentMethod != "cash" && transaction.PaymentMethod != "card" && transaction.PaymentMethod != "transfer"{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "payment_method must be a `cash`, `card` or `transfer`",
		})
		return
	}

	if transaction.CreatedAt.IsZero(){
		transaction.CreatedAt = time.Now()
	}

	if err := h.storage.CreateTransaction(&transaction); err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create transaction: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "transaction created successfully",
		"transaction": transaction,
	})
}

func(h *TransactionHandler) UpdateTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	var transaction models.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	transaction.ID = id

	if transaction.Amount <=0{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "transaction amount must be positive",
		})
		return
	}

	if transaction.Type != "income" && transaction.Type != "expense"{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "transaction type must be 'income' or 'expense'",
		})
		return
	}

	if transaction.PaymentMethod != "cash" && transaction.PaymentMethod != "card" && transaction.PaymentMethod != "transfer"{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "payment_method must be a `cash`, `card` or `transfer`",
		})
		return
	}

	if err := h.storage.UpdateTransaction(&transaction); err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update transaction: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transaction updated successfully",
		"transaction": transaction,
	})
}

func(h *TransactionHandler) DeleteTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	if err := h.storage.DeleteTransaction(id); err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete transaction: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transaction deleted successfully",
	})
}