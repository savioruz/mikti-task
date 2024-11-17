package todo

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
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

func (r *TodoRepositoryImpl) GetPaginated(db *gorm.DB, todos *[]entity.Todo, opts model.TodoQueryOptions) (int64, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.Size <= 0 {
		opts.Size = 10
	}

	query := r.buildPaginatedQuery(db, opts)

	// Get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return 0, err
	}

	// Get paginated results
	offset := (opts.Page - 1) * opts.Size
	if err := query.Offset(offset).Limit(opts.Size).Find(todos).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *TodoRepositoryImpl) buildPaginatedQuery(db *gorm.DB, opts model.TodoQueryOptions) *gorm.DB {
	query := db.Model(&entity.Todo{})

	// Add user filter only if not admin or if admin specified a userID
	if !opts.IsAdmin || (opts.IsAdmin && opts.UserID != nil) {
		query = query.Where("user_id = ?", opts.UserID)
	}

	// Add title filter if provided
	if opts.Title != nil && *opts.Title != "" {
		query = query.Where("title LIKE ?", "%"+*opts.Title+"%")
	}

	return query.Order("created_at DESC")
}
