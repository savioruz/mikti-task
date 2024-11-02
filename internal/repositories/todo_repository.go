package repositories

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoRepository struct {
	Repository[entities.Todo]
	Log *logrus.Logger
}

func NewTodoRepository(log *logrus.Logger) *TodoRepository {
	return &TodoRepository{
		Log: log,
	}
}

func (r *TodoRepository) GetByID(db *gorm.DB, todo *entities.Todo, id string) error {
	return db.Where("id = ?", id).Take(&todo).Error
}

func (r *TodoRepository) GetAll(db *gorm.DB, todos *[]entities.Todo, page, size int) (int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	var totalCount int64
	countResult := db.Model(&entities.Todo{}).Count(&totalCount)
	if countResult.Error != nil {
		return 0, countResult.Error
	}

	result := db.Offset(offset).Limit(size).Find(todos)
	if result.Error != nil {
		return 0, result.Error
	}

	return totalCount, nil
}
