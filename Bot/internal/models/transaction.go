package models

type CreateTransaction struct {
	Category string
	UserID   string
	Name     string
	Cost     float32
}

type Transaction struct {
	ID       string
	UserID   string
	Category string
	Name     string
	Cost     float32
	Date     string
}
