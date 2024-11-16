package repositories

import "gorm.io/gorm"

type RepositoryImpl[T any] struct {
	DB *gorm.DB
}

func (r *RepositoryImpl[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(&entity).Error
}

func (r *RepositoryImpl[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(&entity).Error
}

func (r *RepositoryImpl[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(&entity).Error
}
