package entities

import "gorm.io/gorm"

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Token    string `json:"token" gorm:"not null"`
	gorm.Model
}
