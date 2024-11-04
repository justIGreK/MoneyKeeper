package models

import "time"

type CreateBudget struct{
	UserID    string    `json:"userid" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	EndTime string `json:"endtime" validate:"required"`
}

type Budget struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	Name      string    `bson:"name"`
	Amount    float64   `bson:"amount"`
	DailyAmount float64 `bson:"dailyAmount"`
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsActive  bool      `bson:"is_active"`
}


