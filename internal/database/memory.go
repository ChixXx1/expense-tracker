package database

import (
	"github.com/ChixXx1/expense-tracker/internal/models"
)

type MemoryStorage struct {
	//mu 						sync.Mutex
	//budgets 			[]models.Budget
	transactions []models.Transaction
	categories   []models.Category
	nextID       int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		categories: models.GetDefaultCategories(),
		nextID:     1,
	}
}
