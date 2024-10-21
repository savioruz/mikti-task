package models

import "gorm.io/gorm"

type Book struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Author      string `json:"author" gorm:"not null"`
	Year        int    `json:"year" gorm:"not null"`
	BookStatus  int    `json:"book_status" gorm:"not null"`
	Picture     string `json:"picture,omitempty"`
	Description string `json:"description,omitempty"`
	Rating      int    `json:"rating,omitempty"`
	gorm.Model
}

type BookResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Year        int    `json:"year"`
	BookStatus  int    `json:"book_status"`
	Picture     string `json:"picture,omitempty"`
	Description string `json:"description,omitempty"`
	Rating      int    `json:"rating,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BookCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Author      string `json:"author" validate:"required"`
	Year        int    `json:"year" validate:"required,min=100,max=2200"`
	BookStatus  int    `json:"book_status" validate:"required,len=1"`
	Picture     string `json:"picture,omitempty" validate:"omitempty,url"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Rating      int    `json:"rating,omitempty" validate:"omitempty,numeric,min=1,max=5"`
}

type BookIdRequest struct {
	ID string `json:"id" param:"id" validate:"required,uuid"`
}

type BookUpdateRequest struct {
	Title       string `json:"title,omitempty" validate:"omitempty"`
	Author      string `json:"author,omitempty" validate:"omitempty"`
	Year        int    `json:"year,omitempty" validate:"required,min=100,max=2200"`
	BookStatus  int    `json:"book_status,omitempty" validate:"omitempty,len=1"`
	Picture     string `json:"picture,omitempty" validate:"omitempty,url"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Rating      int    `json:"rating,omitempty" validate:"omitempty,numeric,min=1,max=5"`
}
