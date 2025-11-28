package database

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/ChixXx1/expense-tracker/internal/models"
)

type JSONStorage struct {
	memory   *MemoryStorage
	filepath string
	mu       sync.Mutex
}

func NewJSONStorage(filepath string) *JSONStorage{
	storage := &JSONStorage{
		memory: 	NewMemoryStorage(),
		filepath: filepath,
		mu: 			sync.Mutex{},
	}

	storage.LoadFromFile()

	return storage
}

func(s *JSONStorage) LoadFromFile() error {
	fileContent, err := os.ReadFile(s.filepath)
	if err != nil {
		//log.Fatalf("ошибка чтения файла...")
		return nil
	}

	var data struct {
		Categories []models.Category `json:"categories"`
	}

	if err := json.Unmarshal(fileContent, &data); err != nil {
		return err
	}

	s.memory.categories = data.Categories

	return nil
}

func(s *JSONStorage) SaveToFile() error {
	data := struct {
		Categories []models.Category `json:"categories"`
	}{
		Categories: s.memory.categories,
	}

	file, err := json.MarshalIndent(&data, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, file, 0644)
}

func(s *JSONStorage) GetCategories() ([]models.Category, error){
	return s.memory.categories, nil
}