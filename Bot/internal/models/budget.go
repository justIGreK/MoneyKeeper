package models
type User struct{
	UserID string `bson:"user_id"`
	ChatID string `bson:"chat_id"`
}

type Budget struct{
	ID string 
	Name string
	Amount float32
	DailyAmount float32
	StartDate string
	EndDate string
	CreatedAt string
	UpdatedAt string
	IsActive bool
}

type CreateBudget struct{
	UserID string
	Name string 
	Amount float32
	EndTime string
}