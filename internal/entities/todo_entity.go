package entities

import (
	"time"
)

type Todo struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Title     string    `json:"title" gorm:"not null"`
	Completed bool      `json:"completed" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
