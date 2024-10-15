package entities

import "gorm.io/gorm"

type Todo struct {
	ID        string `json:"id" gorm:"primary_key"`
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"not null"`
	gorm.Model
}
