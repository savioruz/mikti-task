package utils

import "gorm.io/gorm"

type Query[T any] struct {
	DB *gorm.DB
}

func (q *Query[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(&entity).Error
}

func (q *Query[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(&entity).Error
}

func (q *Query[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(&entity).Error
}

func (q *Query[T]) Find(db *gorm.DB, entity *T, id string) error {
	return db.First(&entity, id).Error
}

func (q *Query[T]) FindAll(db *gorm.DB, entity *[]T) error {
	return db.Find(&entity).Error
}

func (q *Query[T]) FindByID(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Find(&entity).Error
}

func (q *Query[T]) FindByField(db *gorm.DB, entity *T, field string, value any) error {
	return db.Where(field+" = ?", value).Find(&entity).Error
}

func (q *Query[T]) CountByID(db *gorm.DB, entity *T, id string) (int64, error) {
	var count int64
	err := db.Model(&entity).Where("id = ?", id).Count(&count).Error
	return count, err
}
