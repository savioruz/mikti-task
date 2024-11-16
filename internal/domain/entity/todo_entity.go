package entity

import "gorm.io/gorm"

type Todo struct {
	ID        string `json:"id" gorm:"primary_key"`
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"not null"`
	UserID    string `json:"user_id" gorm:"not null"`
	User      User   `json:"user" gorm:"foreignKey:UserID"`
	gorm.Model
}
