package models

type ReportResponse struct {
    UserID          string      
    Period          string     
    TotalSpent      float32    
    TransactionCount int32   
    Categories      []*CategoryReport 
}

type CategoryReport struct {
    Name  string  
    Total float32
	Count int32    
}

type BudgetReport struct{
	TotalSpent float32
	BudgetExceeded bool
	FailedDays int32
	SuccessfulDays int32
	MostExpensive Transaction
	MostFrequent string
}