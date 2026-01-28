package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/ChixXx1/expense-tracker/internal/models"
)

type JSONStorage struct {
	mu           sync.RWMutex
	categories   []models.Category
	transactions []models.Transaction
	filepath     string
	nextID       map[string]int
}

func NewJSONStorage(filepath string) *JSONStorage {
	storage := &JSONStorage{
		mu: sync.RWMutex{},
		//categories: 	models.GetDefaultCategories(),
		filepath: filepath,
		nextID: map[string]int{
			"category":    1,
			"transaction": 1,
		},
	}

	if err := storage.load(); err != nil {
		storage.categories = models.GetDefaultCategories()
		storage.transactions = []models.Transaction{}
		storage.updateNextID()
		storage.save()
	}

	return storage
}

func (s *JSONStorage) updateNextID() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	catMaxID := 0
	transMaxID := 0
	for _, cat := range s.categories {
		if cat.ID > catMaxID {
			catMaxID = cat.ID
		}
	}

	for _, tr := range s.transactions {
		if tr.ID > transMaxID {
			transMaxID = tr.ID
		}
	}

	s.nextID["category"] = catMaxID + 1
	s.nextID["transaction"] = transMaxID + 1
}

func (s *JSONStorage) load() error {
	fileData, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	var data struct {
		Categories   []models.Category    `json:"categories"`
		Transactions []models.Transaction `json:"transactions"`
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return err
	}

	s.mu.Lock()
	s.categories = data.Categories
	s.transactions = data.Transactions
	s.mu.Unlock()

	s.updateNextID()

	return nil
}
func (s *JSONStorage) save() error {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

	data := struct {
		Categories   []models.Category    `json:"categories"`
		Transactions []models.Transaction `json:"transactions"`
	}{
		Categories:   s.categories,
		Transactions: s.transactions,
	}

	fileData, err := json.MarshalIndent(&data, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, fileData, 0644)
}

func (s *JSONStorage) GetCategories() ([]models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	categories := make([]models.Category, len(s.categories))
	copy(categories, s.categories)

	return categories, nil
}
func (s *JSONStorage) GetCategoryByID(id int) (*models.Category, error) {
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
func (s *JSONStorage) CreateCategory(category *models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := category.Validate(); err != nil {
		return err
	}

	for _, cat := range s.categories {
		if cat.Name == category.Name && cat.Type == category.Type {
			return errors.New("category with this name is already exists with this type")
		}
	}

	category.ID = s.nextID["category"]
	s.nextID["category"]++
	s.categories = append(s.categories, *category)

	return s.save()
}
func (s *JSONStorage) UpdateCategory(category *models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := category.Validate(); err != nil {
		return err
	}

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
func (s *JSONStorage) DeleteCategory(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, cat := range s.categories {
		if id == cat.ID {
			s.categories = append(s.categories[:i], s.categories[i+1:]...)
			return s.save()
		}
	}

	return errors.New("category is not found")
}

func (s *JSONStorage) GetTransactions(filters TransactionFilters) ([]models.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Transaction

	for _, tr := range s.transactions {
		if filters.StartDate != nil && tr.Date.Before(*filters.StartDate) {
			continue
		}

		if filters.EndDate != nil && tr.Date.After(*filters.EndDate) {
			continue
		}

		if filters.CategoryID != nil && tr.CategoryID != *filters.CategoryID {
			continue
		}

		if filters.Type != nil && tr.Type != *filters.Type {
			continue
		}

		if filters.PaymentMethod != nil && tr.PaymentMethod != *filters.PaymentMethod {
			continue
		}

		result = append(result, tr)
	}

	start := 0
	if filters.Offset != nil && *filters.Offset > 0 {
		start = *filters.Offset
	}

	if start >= len(result) {
		return []models.Transaction{}, nil
	}

	end := len(result)
	if filters.Limit != nil && *filters.Limit > 0 {
		end = start + *filters.Limit
		if end > len(result) {
			end = len(result)
		}
	}

	return result[start:end], nil
}
func (s *JSONStorage) GetTransactionByID(id int) (*models.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, tr := range s.transactions {
		if tr.ID == id {
			transaction := s.transactions[i]
			return &transaction, nil
		}
	}

	return nil, errors.New("transaction not found")
}
func (s *JSONStorage) CreateTransaction(transaction *models.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := transaction.Validate(); err != nil {
		return err
	}

	categoryExists := false
	for _, cat := range s.categories {
		if cat.ID == transaction.CategoryID {
			categoryExists = true
			break
		}
	}

	if !categoryExists {
		return errors.New("category does not exist")
	}

	transaction.ID = s.nextID["transaction"]
	s.nextID["transaction"]++

	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
	}

	s.transactions = append(s.transactions, *transaction)

	return s.save()
}
func (s *JSONStorage) UpdateTransaction(transaction *models.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := transaction.Validate(); err != nil {
		return err
	}

	categoryExists := false
	for _, cat := range s.categories {
		if cat.ID == transaction.CategoryID {
			categoryExists = true
			break
		}
	}

	if !categoryExists {
		return errors.New("category does not exist")
	}

	for i, tr := range s.transactions {
		if tr.ID == transaction.ID {
			s.transactions[i] = *transaction
			return s.save()
		}
	}

	return errors.New("transaction not found")
}
func (s *JSONStorage) DeleteTransaction(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, tr := range s.transactions {
		if id == tr.ID {
			s.transactions = append(s.transactions[:i], s.transactions[i+1:]...)
			return s.save()
		}
	}

	return errors.New("transaction is not found")
}
