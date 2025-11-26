package models

import "time"

type FinancialSummary struct {
	TotalIncome   float64 	`json:"total_income"`
	TotalExpenses float64 	`json:"total_expenses"`
	Balance       float64 	`json:"balance"`
	Period        string 		`json:"period"`
	StartDate     time.Time `json:"start_date"`
	EndDate				time.Time `json:"end_date"`
}

type CategorySummary struct {
	CategoryID 		int 			`json:"category_id"`
	CategoryName  string 		`json:"category_name"`
	Amount 				float64 	`json:"amount"`
	Persentage 		float64 	`json:"persentage"`
	Type 					string 		`json:"type"`
}

type BudgetReport struct {
	Budget 				Budget 		`json:"budget"`
	SpentAmount 	float64 	`json:"spent_amount"`
	Remaining 		float64 	`json:"remaining"`
	Progress 			float64 	`json:"progress"`
	IsOverBudget 	bool 			`json:"is_over_budget"`
}