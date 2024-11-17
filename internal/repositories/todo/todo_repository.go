package todo

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/savioruz/mikti-task/internal/repositories"
	"gorm.io/gorm"
)

type TodoRepository interface {
	repositories.Repository[entity.Todo]
	GetByID(db *gorm.DB, todo *entity.Todo, id string) error
	GetPaginated(db *gorm.DB, todos *[]entity.Todo, opts model.TodoQueryOptions) (int64, error)
}
