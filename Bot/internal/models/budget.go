package models
type User struct{
	UserID string `bson:"user_id"`
	ChatID string `bson:"chat_id"`
}

type Budget struct{
	ID string 
	Name string
	Limit float32
	Start string
	End string
	Category  []Category 
}

type CreateBudget struct{
	UserID string
	Name string 
	Amount float32
	Period string
	Start string
	End string
}

type Category struct {
	ID    string  
	Name  string  
	Limit float32 
}

type CreateCategory struct{
	UserID string
	BudgetID string
	Name string
	Limit float32
}

type UpdateBudget struct{
	BudgetID string 
	UserID string 
	Name *string 
	Limit *float64 
	Start *string
	End *string
}
type UpdateCategory struct{
	UserID string
	BudgetID string
	CategoryID string
	Name *string
	Limit *float64
}
