package database

import "github.com/ChixXx1/expense-tracker/internal/models"

type Storage interface {
	GetCategories() ([]models.Category, error)
	GetCategoryByID(id int) (*models.Category, error)
	CreateCategory() error
	UpdateCategory() error 
	DeleteCategory()
}