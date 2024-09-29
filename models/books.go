package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title           string `json:"title"`
	Author          string `json:"author"`
	Copies          int    `json:"copies"`
	AvailableCopies int    `json:"available_copies"`
}

type Borrower struct {
	gorm.Model
	BookID uint   `json:"book_id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "borrowed" or "returned"
}
