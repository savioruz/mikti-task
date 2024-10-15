package usecases

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/tree/week-3/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models/converter"
	"github.com/savioruz/mikti-task/tree/week-3/internal/repositories"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TodoUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	TodoRepository *repositories.TodoRepository
}

func NewTodoUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, todoRepository *repositories.TodoRepository) *TodoUsecase {
	return &TodoUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		TodoRepository: todoRepository,
	}
}

func (u *TodoUsecase) Create(ctx context.Context, request *models.TodoCreateRequest) (*models.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	todo := &entities.Todo{
		ID:        uuid.New().String(),
		Title:     request.Title,
		Completed: false,
	}

	if err := u.TodoRepository.Create(tx, todo); err != nil {
		u.Log.Errorf("failed to create todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.TodoToResponse(todo), nil
}

func (u *TodoUsecase) Update(ctx context.Context, id *models.TodoUpdateIDRequest, request *models.TodoUpdateRequest) (*models.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.Completed == nil && request.Title == nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	if err := u.Validate.Struct(id); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	todo := &entities.Todo{}
	if err := u.TodoRepository.GetByID(tx, todo, id.ID); err != nil {
		u.Log.Errorf("failed to get todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if request != nil {
		if request.Title != nil {
			todo.Title = *request.Title
		}
		if request.Completed != nil {
			todo.Completed = *request.Completed
		}
	}
	todo.UpdatedAt = time.Now()

	if err := u.TodoRepository.Update(tx, todo); err != nil {
		u.Log.Errorf("failed to update todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.TodoToResponse(todo), nil
}

func (u *TodoUsecase) Delete(ctx context.Context, request *models.TodoDeleteRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	todo := &entities.Todo{}
	if err := u.TodoRepository.GetByID(tx, todo, request.ID); err != nil {
		u.Log.Errorf("failed to get todo: %v", err)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := u.TodoRepository.Delete(tx, todo); err != nil {
		u.Log.Errorf("failed to delete todo: %v", err)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (u *TodoUsecase) Get(ctx context.Context, request *models.TodoGetRequest) (*models.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	todo := &entities.Todo{}
	if err := u.TodoRepository.GetByID(tx, todo, request.ID); err != nil {
		u.Log.Errorf("failed to get todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.TodoToResponse(todo), nil
}

func (u *TodoUsecase) List(ctx context.Context, request *models.TodoListRequest) ([]*models.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	var todos []entities.Todo
	if err := u.TodoRepository.GetAll(tx, &todos); err != nil {
		u.Log.Errorf("failed to get todos: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	var todoResponses []*models.TodoResponse
	for _, todo := range todos {
		todoResponses = append(todoResponses, converter.TodoToResponse(&todo))
	}

	return todoResponses, nil
}
