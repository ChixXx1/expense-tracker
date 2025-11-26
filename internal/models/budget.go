package models

import "time"

const (
	BudgetPeriodMonthly = "monthly"
	BudgetPeriodWeekly = "weekly"
	BudgetPeriodYearly = "yearly"
)

type Budget struct {
	ID         int 				`json:"id"`
	CategoryID int 				`json:"category_id"`
	Amount     float64 		`json:"amount"`
	Period     string 		`json:"period"`
	Month      time.Time  `json:"month"`
	Spent      float64 		`json:"spent"`
	CreatedAt  time.Time 	`json:"created_at"`
}