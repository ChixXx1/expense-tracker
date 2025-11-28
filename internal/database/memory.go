package database

import (
	"github.com/ChixXx1/expense-tracker/internal/models"
)

type MemoryStorage struct {
	//mu 						sync.Mutex
	//transactions 	[]models.Transaction
	//budgets 			[]models.Budget
	categories 		[]models.Category
	nextID 				int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		categories: 	models.GetDefaultCategories(),
		nextID: 			1,
	}
}

func(ms *MemoryStorage) GetCategories() ([]models.Category, error) {
	return ms.categories, nil
}