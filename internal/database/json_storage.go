package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ChixXx1/expense-tracker/internal/models"
)

type JSONStorage struct {
	mu       			sync.RWMutex
	categories 		[]models.Category
	filepath 			string
	nextID	 			map[string]int
}

func NewJSONStorage(filepath string) *JSONStorage{
	storage := &JSONStorage{
		mu: 					sync.RWMutex{},
		//categories: 	models.GetDefaultCategories(),
		filepath: 		filepath,
		nextID: 			map[string]int{
			"category": 1,
		},
	}

	if err := storage.load(); err != nil{
		storage.categories = models.GetDefaultCategories()

		storage.save()
	}

	return storage
}

func(s *JSONStorage) load() error {
	fileData, err := os.ReadFile(s.filepath)
	if err != nil{
		return err
	}

	var data struct{
		Categories []models.Category `json:"categories"`
	}
	
	if err := json.Unmarshal(fileData, &data); err != nil {
		return err
	}

	s.mu.Lock()
	s.categories = data.Categories
	s.mu.Unlock()

	return nil
}

func(s *JSONStorage) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := struct{
		Categories []models.Category `json:"categories"`
	}{
		Categories: s.categories,
	}

	fileData, err := json.MarshalIndent(&data, "", "	")
	if err != nil {
		return err
	}
	
	return os.WriteFile(s.filepath, fileData, 0644)
}

func(s *JSONStorage) GetCategories() ([]models.Category, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	categories := make([]models.Category, len(s.categories))
	copy(categories, s.categories)

	return categories, nil
}

func(s *JSONStorage) GetCategoryByID(id int) (*models.Category, error){
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.categories {
		if id == s.categories[i].ID {
			category := &s.categories[i]
			return category, nil
		}
	}

	return nil, errors.New("category is not found!")
}

func(s *JSONStorage) CreateCategory(category *models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, cat := range s.categories {
		if cat.Name == category.Name && cat.Type == category.Type{
			return errors.New("category with this name is already exists with this type")
		}
	}
	return nil
}