package models

import (
	"errors"
	"time"
)

const (
	BudgetPeriodMonthly = "monthly"
	BudgetPeriodWeekly  = "weekly"
	BudgetPeriodYearly  = "yearly"
)

type Budget struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Amount     float64   `json:"amount"`
	Period     string    `json:"period"`
	Month      time.Time `json:"month"`
	Spent      float64   `json:"spent"`
	CreatedAt  time.Time `json:"created_at"`
}

func (b *Budget) Validate() error {
	if b.Amount <= 0 {
		return errors.New("budget amount must be positive")
	}

	if b.CategoryID <= 0 {
		return errors.New("category_id must be positive")
	}

	validPeriods := map[string]bool{
		BudgetPeriodMonthly: true,
		BudgetPeriodWeekly:  true,
		BudgetPeriodYearly:  true,
	}

	if !validPeriods[b.Period] {
		return errors.New("period must be 'monthly', 'weekly' or 'yearly'")
	}

	if b.Month.IsZero() {
		return errors.New("month is required")
	}

	return nil
}
