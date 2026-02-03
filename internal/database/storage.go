package database

import (
	"time"

	"github.com/ChixXx1/expense-tracker/internal/models"
)

type BudgetFilters struct {
	CategoryID *int
	Period     *string
	Month      *time.Time
}

type TransactionFilters struct {
	StartDate     *time.Time
	EndDate       *time.Time
	CategoryID    *int
	Type          *string
	PaymentMethod *string
	Limit         *int
	Offset        *int
}

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

	GetBudgets(filters BudgetFilters) ([]models.Budget, error)
	GetBudgetByID(id int) (*models.Budget, error)
	CreateBudget(budget *models.Budget) error
	UpdateBudget(budget *models.Budget) error
	DeleteBudget(id int) error

	GetFinancialSummary(startDate, endDate time.Time) (*models.FinancialSummary, error)
	GetCategorySummary(startDate, endDate time.Time) ([]models.CategorySummary, error)
	GetBudgetReport(budgetID int) (*models.BudgetReport, error)
}
