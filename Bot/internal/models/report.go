package models

type ReportResponse struct {   
    Period          string     
    TotalSpent      float32    
    TransactionCount int32   
    Categories      []CategoryReport 
}

type CategoryReport struct {
    Name  string  
    Total float32
	Count int32    
}

type RequiredCategoryReport struct {
	Name  string
	Total float32
    Limit float32
	Count int32
}
type BudgetReport struct {
	BudgetName string
	Period     string
	LeftDays   float32
	Limit      float32
	TotalSpent float32
	TransactionCount int
    RequiredCategories []RequiredCategoryReport
	Categories       []CategoryReport
}

type DatesSummary struct{
	UserID string
	Period string
	Start string
	End string
}