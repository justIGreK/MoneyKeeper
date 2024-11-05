package models 

type ReportResponse struct {
    UserID          string      `json:"user_id"`
    Period          string      `json:"period"`
    TotalSpent      float64     `json:"total_spent"`
    TransactionCount int        `json:"transaction_count"`
    Categories      []CategoryReport `json:"categories"`
}

type CategoryReport struct {
    Name  string  `json:"name"`
    Total float64 `json:"total"`
	Count int    `json:"transaction_count"`
}


type BudgetReport struct{
	TotalSpent float64
	BudgetExceeded bool
	FailedDays int
	SuccessfulDays int
	MostExpensive Transaction
	MostFrequent string
}