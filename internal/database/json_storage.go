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
	budgets      []models.Budget
	filepath     string
	nextID       map[string]int
}

func NewJSONStorage(filepath string) *JSONStorage {
	storage := &JSONStorage{
		mu: sync.RWMutex{},
		//categories: 	models.GetDefaultCategories(),
		budgets:  []models.Budget{},
		filepath: filepath,
		nextID: map[string]int{
			"category":    1,
			"transaction": 1,
			"budget":      1,
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
	budgetMaxID := 0

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

	for _, bud := range s.budgets {
		if bud.ID > budgetMaxID {
			budgetMaxID = bud.ID
		}
	}

	s.nextID["category"] = catMaxID + 1
	s.nextID["transaction"] = transMaxID + 1
	s.nextID["budget"] = budgetMaxID + 1
}

func (s *JSONStorage) load() error {
	fileData, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	var data struct {
		Categories   []models.Category    `json:"categories"`
		Transactions []models.Transaction `json:"transactions"`
		Budgets      []models.Budget      `json:"budgets"`
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return err
	}

	s.mu.Lock()
	s.categories = data.Categories
	s.transactions = data.Transactions
	s.budgets = data.Budgets
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
		Budgets      []models.Budget      `json:"budgets"`
	}{
		Categories:   s.categories,
		Transactions: s.transactions,
		Budgets:      s.budgets,
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

func (s *JSONStorage) GetBudgets(filters BudgetFilters) ([]models.Budget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Budget

	for _, budget := range s.budgets {
		if filters.CategoryID != nil && budget.CategoryID != *filters.CategoryID {
			continue
		}

		if filters.Period != nil && budget.Period != *filters.Period {
			continue
		}

		if filters.Month != nil {
			// Сравниваем только год и месяц
			yearMatch := budget.Month.Year() == filters.Month.Year()
			monthMatch := budget.Month.Month() == filters.Month.Month()
			if !(yearMatch && monthMatch) {
				continue
			}
		}

		result = append(result, budget)
	}

	return result, nil
}

func (s *JSONStorage) GetBudgetByID(id int) (*models.Budget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, budget := range s.budgets {
		if budget.ID == id {
			return &s.budgets[i], nil
		}
	}

	return nil, errors.New("budget not found")
}

func (s *JSONStorage) CreateBudget(budget *models.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Валидация
	if err := budget.Validate(); err != nil {
		return err
	}

	// Проверка существования категории
	categoryExists := false
	for _, cat := range s.categories {
		if cat.ID == budget.CategoryID {
			categoryExists = true
			break
		}
	}

	if !categoryExists {
		return errors.New("category does not exist")
	}

	// Проверка дубликата бюджета для категории в тот же период
	for _, existing := range s.budgets {
		if existing.CategoryID == budget.CategoryID &&
			existing.Period == budget.Period &&
			existing.Month.Year() == budget.Month.Year() &&
			existing.Month.Month() == budget.Month.Month() {
			return errors.New("budget already exists for this category and period")
		}
	}

	budget.ID = s.nextID["budget"]
	s.nextID["budget"]++

	if budget.CreatedAt.IsZero() {
		budget.CreatedAt = time.Now()
	}

	s.budgets = append(s.budgets, *budget)

	return s.save()
}

func (s *JSONStorage) UpdateBudget(budget *models.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := budget.Validate(); err != nil {
		return err
	}

	for i, existing := range s.budgets {
		if existing.ID == budget.ID {
			s.budgets[i] = *budget
			return s.save()
		}
	}

	return errors.New("budget not found")
}

func (s *JSONStorage) DeleteBudget(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, budget := range s.budgets {
		if budget.ID == id {
			s.budgets = append(s.budgets[:i], s.budgets[i+1:]...)
			return s.save()
		}
	}

	return errors.New("budget not found")
}

func (s *JSONStorage) GetFinancialSummary(startDate, endDate time.Time) (*models.FinancialSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalIncome, totalExpenses float64

	for _, tx := range s.transactions {
		if tx.Date.Before(startDate) || tx.Date.After(endDate) {
			continue
		}

		// Используем switch вместо if-else (рекомендация staticcheck)
		switch tx.Type {
		case models.TransactionTypeIncome:
			totalIncome += tx.Amount
		case models.TransactionTypeExpense:
			totalExpenses += tx.Amount
			// default: игнорируем неизвестные типы
		}
	}

	balance := totalIncome - totalExpenses

	// Определяем период
	period := "custom"
	if startDate.Year() == endDate.Year() && startDate.Month() == endDate.Month() {
		period = "monthly"
	}

	return &models.FinancialSummary{
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		Balance:       balance,
		Period:        period,
		StartDate:     startDate,
		EndDate:       endDate,
	}, nil
}

func (s *JSONStorage) GetCategorySummary(startDate, endDate time.Time) ([]models.CategorySummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categoryAmounts := make(map[int]float64)
	categoryTypes := make(map[int]string)
	categoryNames := make(map[int]string)

	// Собираем суммы по категориям
	for _, tx := range s.transactions {
		if tx.Date.Before(startDate) || tx.Date.After(endDate) {
			continue
		}

		categoryAmounts[tx.CategoryID] += tx.Amount
	}

	// Получаем имена и типы категорий
	for _, cat := range s.categories {
		if _, exists := categoryAmounts[cat.ID]; exists {
			categoryTypes[cat.ID] = cat.Type
			categoryNames[cat.ID] = cat.Name
		}
	}

	// Считаем общий итог для расчета процентов
	var totalAmount float64
	for _, amount := range categoryAmounts {
		totalAmount += amount
	}

	// Формируем результат
	var summaries []models.CategorySummary
	for categoryID, amount := range categoryAmounts {
		percentage := 0.0
		if totalAmount > 0 {
			percentage = (amount / totalAmount) * 100
		}

		summaries = append(summaries, models.CategorySummary{
			CategoryID:   categoryID,
			CategoryName: categoryNames[categoryID],
			Amount:       amount,
			Persentage:   percentage, // Опечатка в модели, но оставляем как есть
			Type:         categoryTypes[categoryID],
		})
	}

	return summaries, nil
}

func (s *JSONStorage) GetBudgetReport(budgetID int) (*models.BudgetReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Находим бюджет
	var budget *models.Budget
	for i := range s.budgets {
		if s.budgets[i].ID == budgetID {
			budget = &s.budgets[i]
			break
		}
	}

	if budget == nil {
		return nil, errors.New("budget not found")
	}

	// Находим категорию бюджета
	/* var categoryName string
	for _, cat := range s.categories {
		if cat.ID == budget.CategoryID {
			categoryName = cat.Name
			break
		}
	} */

	// Считаем потраченную сумму за период
	var spentAmount float64

	// Определяем период для фильтрации транзакций
	startDate := budget.Month
	endDate := budget.Month

	switch budget.Period {
	case models.BudgetPeriodMonthly:
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case models.BudgetPeriodWeekly:
		endDate = startDate.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case models.BudgetPeriodYearly:
		endDate = startDate.AddDate(1, 0, 0).Add(-time.Nanosecond)
	}

	for _, tx := range s.transactions {
		if tx.CategoryID == budget.CategoryID &&
			!tx.Date.Before(startDate) &&
			!tx.Date.After(endDate) {
			spentAmount += tx.Amount
		}
	}

	remaining := budget.Amount - spentAmount
	progress := 0.0
	if budget.Amount > 0 {
		progress = (spentAmount / budget.Amount) * 100
	}

	return &models.BudgetReport{
		Budget:       *budget,
		SpentAmount:  spentAmount,
		Remaining:    remaining,
		Progress:     progress,
		IsOverBudget: spentAmount > budget.Amount,
	}, nil
}
