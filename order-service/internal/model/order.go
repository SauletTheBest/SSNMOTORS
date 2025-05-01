package model

type Product struct {
	ProductID string
	Quantity  int
}

type Order struct {
	ID        string
	UserID    string
	Products  []Product
	Total     float64
	Status    string
}