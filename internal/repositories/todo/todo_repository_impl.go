package todo

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/repositories"
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

func (r *TodoRepositoryImpl) GetByTitle(db *gorm.DB, todo *[]entity.Todo, title, userID string, page, size int) (int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	query := db.Model(&entity.Todo{}).Where("title LIKE ? AND user_id = ?", "%"+title+"%", userID)

	var totalCount int64
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return 0, countResult.Error
	}

	result := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&todo)
	if result.Error != nil {
		return 0, result.Error
	}

	return totalCount, nil
}

func (r *TodoRepositoryImpl) GetAll(db *gorm.DB, todos *[]entity.Todo, userID string, page, size int) (int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	query := db.Model(&entity.Todo{}).Where("user_id = ?", userID)

	var totalCount int64
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return 0, countResult.Error
	}

	result := query.Order("created_at DESC").Offset(offset).Limit(size).Find(todos)
	if result.Error != nil {
		return 0, result.Error
	}

	return totalCount, nil
}
