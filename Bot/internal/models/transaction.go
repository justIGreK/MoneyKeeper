package models

type CreateTransaction struct {
	Category string
	UserID   string
	Name     string
	Cost     float32
	Date *string
}

type Transaction struct {
	ID       string
	UserID   string
	Category string
	Name     string
	Cost     float32
	Date     string
}

type UpdateTransaction struct {
	ID       string
	UserID   string
	Category *string
	Name     *string
	Cost     *float64
	Date     *string
	Time     *string
}
