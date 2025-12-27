package database

import (
	"time"

	"github.com/ChixXx1/expense-tracker/internal/models"
)

type Storage interface {
	GetCategories() ([]models.Category, error)
	GetCategoryByID(id int) (*models.Category, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error 
	DeleteCategory(id int) error

	GetTransactions(filters TransactionFilters) ([]models.Transaction, error)
	GetTransactionByID(id int) (*models.Transaction, error)
	CreateTransaction(transaction *models.Transaction) error
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id int) error 
 
	
}

type TransactionFilters struct{
	StartDate *time.Time
	EndDate *time.Time
	CategoryID *int
	Type *string
	PaymentMethod *string
	Limit *int
	Offset *int
}