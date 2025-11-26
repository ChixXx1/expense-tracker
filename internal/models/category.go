package models

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
	ParentID *int   `json:"parent_id,omitempty"` //(указатель на int, так как может быть nil для корневых категорий)
}

func GetDefaultCategories() []Category {
	return []Category{
		// Категории расходов
		{ID: 1, Name: "Еда", Type: "expense", Color: "#FF6B6B", Icon: "🍕"},
		{ID: 2, Name: "Транспорт", Type: "expense", Color: "#4ECDC4", Icon: "🚗"},
		{ID: 3, Name: "Развлечения", Type: "expense", Color: "#45B7D1", Icon: "🎬"},
		{ID: 4, Name: "Одежда", Type: "expense", Color: "#FFEAA7", Icon: "👕"},

		// Категории доходов
		{ID: 5, Name: "Зарплата", Type: "income", Color: "#A8E6CF", Icon: "💰"},
		{ID: 6, Name: "Фриланс", Type: "income", Color: "#DCEDC1", Icon: "💻"},
		{ID: 7, Name: "Инвестиции", Type: "income", Color: "#FFD3B6", Icon: "📈"},
	}
}

func (c *Category) IsIncomeCategory() bool {
	return c.Type == "income"
}

func (c *Category) IsExpenseCategory() bool {
	return c.Type == "expense"
}