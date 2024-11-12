package repositories

import "gorm.io/gorm"

type Repository[T any] interface {
	Create(db *gorm.DB, entity *T) error
	Update(db *gorm.DB, entity *T) error
	Delete(db *gorm.DB, entity *T) error
	FindByID(db *gorm.DB, entity *T, id any) error
	CountByID(db *gorm.DB, id any) (int64, error)
}
