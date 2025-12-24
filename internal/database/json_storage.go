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
	transactions  []models.Transaction
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
			"transaction": 1,
		},
	}

	if err := storage.load(); err != nil{
		storage.categories = models.GetDefaultCategories()
		storage.updateNextID()
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
	
	s.updateNextID()
	
	return nil
}

func(s *JSONStorage) save() error {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

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

func(s *JSONStorage) updateNextID(){
	s.mu.RLock()
	defer s.mu.RUnlock()

	maxID := 0
	for _, cat := range s.categories{
		if cat.ID > maxID{
			maxID = cat.ID
		}
	}
	s.nextID["category"] = maxID + 1
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

	category.ID = s.nextID["category"]
	s.nextID["category"]++
	s.categories = append(s.categories, *category)

	return s.save()
}

func(s *JSONStorage) UpdateCategory(category *models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, cat := range s.categories {
		if cat.ID == category.ID {
			for j, other := range s.categories {
				if i != j && category.Name == other.Name && category.Type == other.Type {
					return errors.New("category with this name already exists for this type")
				}
			}
			s.categories[i] = *category
			return s.save()
		}
	}

	return errors.New("category is not found")
}

func(s *JSONStorage) DeleteCategory(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, cat := range s.categories {
		if id == cat.ID {
			s.categories = append(s.categories[:i], s.categories[i + 1:]...)
			return s.save()
		}
	}

	return errors.New("category is not found")
}