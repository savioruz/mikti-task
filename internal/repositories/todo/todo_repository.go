package todo

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/repositories"
	"gorm.io/gorm"
)

type TodoRepository interface {
	repositories.Repository[entity.Todo]
	GetByID(db *gorm.DB, todo *entity.Todo, id string) error
	GetAll(db *gorm.DB, todos *[]entity.Todo, userID string, page, size int) (int64, error)
}
