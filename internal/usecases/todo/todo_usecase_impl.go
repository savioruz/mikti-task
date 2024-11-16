package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/savioruz/mikti-task/internal/domain/model/converter"
	"github.com/savioruz/mikti-task/internal/platform/cache"
	"github.com/savioruz/mikti-task/internal/platform/jwt"
	"github.com/savioruz/mikti-task/internal/repositories/todo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TodoUsecaseImpl struct {
	DB             *gorm.DB
	Cache          *cache.ImplCache
	Log            *logrus.Logger
	Validate       *validator.Validate
	TodoRepository todo.TodoRepository
}

func NewTodoUsecaseImpl(db *gorm.DB, c *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, todoRepository todo.TodoRepository) *TodoUsecaseImpl {
	return &TodoUsecaseImpl{
		DB:             db,
		Cache:          c,
		Log:            log,
		Validate:       validate,
		TodoRepository: todoRepository,
	}
}

func (u *TodoUsecaseImpl) Create(ctx context.Context, request *model.TodoCreateRequest) (*model.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.getJWTClaims(ctx)
	if err != nil {
		u.Log.Errorf("failed to get JWT claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	todoData := &entity.Todo{
		ID:     uuid.NewString(),
		Title:  request.Title,
		Done:   false,
		UserID: claims.UserID,
	}

	if err := u.TodoRepository.Create(tx, todoData); err != nil {
		u.Log.Errorf("failed to create todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	u.invalidateUserListCache(claims.UserID)

	return converter.TodoToResponse(todoData), nil
}

func (u *TodoUsecaseImpl) Update(ctx context.Context, id *model.TodoUpdateIDRequest, request *model.TodoUpdateRequest) (*model.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.Done == nil && request.Title == nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	if err := u.Validate.Struct(id); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	todoData := &entity.Todo{}
	if err := u.TodoRepository.GetByID(tx, todoData, id.ID); err != nil {
		u.Log.Errorf("failed to get todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if request.Title != nil {
		todoData.Title = *request.Title
	}
	if request.Done != nil {
		todoData.Done = *request.Done
	}

	if err := u.TodoRepository.Update(tx, todoData); err != nil {
		u.Log.Errorf("failed to update todo: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.TodoToResponse(todoData), nil
}

func (u *TodoUsecaseImpl) Delete(ctx context.Context, request *model.TodoDeleteRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return false, errors.New(http.StatusText(http.StatusBadRequest))
	}

	todoData := &entity.Todo{}
	if err := u.TodoRepository.GetByID(tx, todoData, request.ID); err != nil {
		u.Log.Errorf("failed to get todo: %v", err)
		return false, errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := u.TodoRepository.Delete(tx, todoData); err != nil {
		u.Log.Errorf("failed to delete todo: %v", err)
		return false, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return false, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return true, nil
}

func (u *TodoUsecaseImpl) Get(ctx context.Context, request *model.TodoGetRequest) (*model.TodoResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	key := fmt.Sprintf("todos:get:%s", request.ID)
	var data *model.TodoResponse
	err := u.Cache.Get(key, &data)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		u.Log.Errorf("failed to get data from cache: %v", err)
	}

	if data != nil {
		u.Log.Infof("data from cache: %v", data)
		return data, nil
	} else {
		tx := u.DB.WithContext(ctx).Begin()
		defer tx.Rollback()

		todoData := &entity.Todo{}
		if err := u.TodoRepository.GetByID(tx, todoData, request.ID); err != nil {
			u.Log.Errorf("failed to get todo: %v", err)
			return nil, errors.New(http.StatusText(http.StatusNotFound))
		}

		response := converter.TodoToResponse(todoData)

		if err := u.Cache.Set(key, response, 5*time.Minute); err != nil {
			u.Log.Errorf("failed to set data to cache: %v", err)
		}

		return response, nil
	}
}

func (u *TodoUsecaseImpl) GetAll(ctx context.Context, request *model.TodoGetAllRequest) (*model.Response[[]*model.TodoResponse], error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.getJWTClaims(ctx)
	if err != nil {
		u.Log.Errorf("failed to get JWT claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	// Ensure valid pagination parameters
	if request.Size <= 0 {
		request.Size = 10 // Default page size
	}
	if request.Page <= 0 {
		request.Page = 1 // Default page number
	}

	// Cache keys for both data and metadata
	dataKey := fmt.Sprintf("todos:user:%s:list:data:%d:%d", claims.UserID, request.Page, request.Size)
	metadataKey := fmt.Sprintf("todos:user:%s:list:metadata", claims.UserID)

	// Try to get cached data
	var cachedData []*model.TodoResponse
	err = u.Cache.Get(dataKey, &cachedData)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		u.Log.Errorf("failed to get data from cache: %v", err)
	}

	// Try to get cached metadata
	var totalItems int
	err = u.Cache.Get(metadataKey, &totalItems)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		u.Log.Errorf("failed to get metadata from cache: %v", err)
	}

	// If we have both cached data and metadata, return them
	if len(cachedData) > 0 && totalItems > 0 {
		u.Log.Infof("Data retrieved from cache")
		totalPages := (totalItems + request.Size - 1) / request.Size
		response := model.NewResponse(cachedData, &model.PageMetadata{
			Page:       request.Page,
			Size:       request.Size,
			TotalItems: totalItems,
			TotalPages: totalPages,
		})

		return response, nil
	}

	// If cache miss, get data from database
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var todos []entity.Todo
	dbTotalItems, err := u.TodoRepository.GetAll(tx, &todos, claims.UserID, request.Page, request.Size)
	if err != nil {
		u.Log.Errorf("failed to get todos from database: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if dbTotalItems == 0 || len(todos) == 0 {
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	todoResponses := converter.TodosToResponses(todos)

	// Cache both the data for this page and the total count
	if err := u.Cache.Set(dataKey, todoResponses, 5*time.Minute); err != nil {
		u.Log.Errorf("failed to set data to cache: %v", err)
	}
	if err := u.Cache.Set(metadataKey, int(dbTotalItems), 5*time.Minute); err != nil {
		u.Log.Errorf("failed to set metadata to cache: %v", err)
	}

	totalPages := (int(dbTotalItems) + request.Size - 1) / request.Size
	response := model.NewResponse(todoResponses, &model.PageMetadata{
		Page:       request.Page,
		Size:       request.Size,
		TotalItems: int(dbTotalItems),
		TotalPages: totalPages,
	})

	return response, nil
}

func (u *TodoUsecaseImpl) invalidateUserListCache(userID string) {
	if err := u.Cache.Delete(fmt.Sprintf("todos:user:%s:list:metadata", userID)); err != nil {
		u.Log.Errorf("failed to delete user metadata cache: %v", err)
	}

	pattern := fmt.Sprintf("todos:user:%s:list:data:*", userID)
	if err := u.Cache.DeletePattern(pattern); err != nil {
		u.Log.Errorf("failed to delete user data caches: %v", err)
	}
}

func (u *TodoUsecaseImpl) getJWTClaims(ctx context.Context) (*jwt.JWTClaims, error) {
	claims, ok := ctx.Value("claims").(*jwt.JWTClaims)
	if !ok || claims == nil {
		return nil, errors.New("unauthorized: invalid or missing JWT claims")
	}
	return claims, nil
}
