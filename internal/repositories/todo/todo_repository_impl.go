package todo

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/domain/entity"
	"github.com/savioruz/mikti-task/tree/week-4/internal/repositories"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoRepositoryImpl struct {
	repositories.RepositoryImpl[entity.Todo]
	Log *logrus.Logger
}

func NewTodoRepository(db *gorm.DB, log *logrus.Logger) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{
		RepositoryImpl: repositories.RepositoryImpl[entity.Todo]{DB: db},
		Log:            log,
	}
}

func (r *TodoRepositoryImpl) GetByID(db *gorm.DB, todo *entity.Todo, id string) error {
	return db.Where("id = ?", id).Take(&todo).Error
}

func (r *TodoRepositoryImpl) GetAll(db *gorm.DB, todos *[]entity.Todo, page, size int) (int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	var totalCount int64
	countResult := db.Model(&entity.Todo{}).Count(&totalCount)
	if countResult.Error != nil {
		return 0, countResult.Error
	}

	result := db.Offset(offset).Limit(size).Find(todos)
	if result.Error != nil {
		return 0, result.Error
	}

	return totalCount, nil
}
