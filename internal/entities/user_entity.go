package entities

import "gorm.io/gorm"

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
	Status   int    `json:"status" gorm:"not null"`
	gorm.Model
}
