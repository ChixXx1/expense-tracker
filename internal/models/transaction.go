package models

import (
	"math"
	"time"
)

//использовать типизированные константы
// type TransactionType string и так далее

const (
	TransactionTypeIncome = "income"
	TransactionTypeExpense = "expense"
)

const (
	PaymentMethodCash = "cash"
	PaymentMethodCard = "card"
	PaymentMethodTransfer = "transfer"
)

type Transaction struct {
	ID         		int 				`json:"id"`
	Amount     		float64 		`json:"amount"`
	Type       		string 			`json:"type"`
	CategoryID 		int 				`json:"category_id"`
	Date       		time.Time 	`json:"date"`
	Description 	string 			`json:"description"`
	PaymentMethod string 			`json:"payment_method"`
	CreatedAt 		time.Time 	`json:"created_at"`
}

func Validate() {}

func(t *Transaction) IsValidAmount() bool {
	return t.Amount > 0 && !math.IsNaN(t.Amount) && !math.IsInf(t.Amount, 0)
}