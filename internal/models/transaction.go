package models

import (
	"errors"
	"math"
	"time"
)

//использовать типизированные константы
// type TransactionType string и так далее

const (
	TransactionTypeIncome  = "income"
	TransactionTypeExpense = "expense"
)

const (
	PaymentMethodCash     = "cash"
	PaymentMethodCard     = "card"
	PaymentMethodTransfer = "transfer"
)

type Transaction struct {
	ID            int       `json:"id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	CategoryID    int       `json:"category_id"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}

	if math.IsNaN(t.Amount) || math.IsInf(t.Amount, 0) {
		return errors.New("transaction amount is invalid")
	}

	if t.Type != TransactionTypeExpense && t.Type != TransactionTypeIncome {
		return errors.New("transaction type must be 'income' or 'expense'")
	}

	if t.CategoryID <= 0 {
		return errors.New("category_id must be positive")
	}

	if t.PaymentMethod != PaymentMethodCash &&
		t.PaymentMethod != PaymentMethodCard &&
		t.PaymentMethod != PaymentMethodTransfer {
		return errors.New("transaction method must be 'cash', 'card' or 'transfer'")
	}

	if t.Date.IsZero() {
		return errors.New("transaction date is required")
	}

	if t.Date.After(time.Now().Add(24 * time.Hour)) {
		return errors.New("transaction date cannot be too far in the future")
	}

	return nil
}

func (t *Transaction) IsValidAmount() bool {
	return t.Amount > 0 && !math.IsNaN(t.Amount) && !math.IsInf(t.Amount, 0)
}
