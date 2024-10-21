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

func (r *TodoRepository) GetAll(db *gorm.DB, todos *[]entities.Todo) error {
	return db.Find(&todos).Error
}