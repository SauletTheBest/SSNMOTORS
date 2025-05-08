package model

import "time"

type Order struct {
	ID        string
	UserID    string
	CarID     string
	Quantity  int32
	Status    string
	CreatedAt time.Time
}
