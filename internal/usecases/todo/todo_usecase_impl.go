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
	"github.com/savioruz/mikti-task/internal/platform/helper"
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
	helper         *helper.ContextHelper
}

func NewTodoUsecaseImpl(db *gorm.DB, c *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, todoRepository todo.TodoRepository) *TodoUsecaseImpl {
	return &TodoUsecaseImpl{
		DB:             db,
		Cache:          c,
		Log:            log,
		Validate:       validate,
		TodoRepository: todoRepository,
		helper:         helper.NewContextHelper(),
	}
}

func (u *TodoUsecaseImpl) Create(ctx context.Context, request *model.TodoCreateRequest) (*model.TodoResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.helper.GetJWTClaims(ctx)
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

	return converter.TodoToResponse(todoData, false), nil
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

	if err := u.helper.VerifyOwnership(ctx, todoData.UserID); err != nil {
		u.Log.Errorf("unauthorized access attempt: %v", err)
		return nil, errors.New(http.StatusText(http.StatusForbidden))
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

	return converter.TodoToResponse(todoData, false), nil
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

	if err := u.helper.VerifyOwnership(ctx, todoData.UserID); err != nil {
		u.Log.Errorf("unauthorized access attempt: %v", err)
		return false, errors.New(http.StatusText(http.StatusForbidden))
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

	isAdmin := u.helper.IsAdmin(ctx)

	if data == nil {
		tx := u.DB.WithContext(ctx).Begin()
		defer tx.Rollback()

		todoData := &entity.Todo{}
		if err := u.TodoRepository.GetByID(tx, todoData, request.ID); err != nil {
			u.Log.Errorf("failed to get todo: %v", err)
			return nil, errors.New(http.StatusText(http.StatusNotFound))
		}

		if err := u.helper.VerifyOwnership(ctx, todoData.UserID); err != nil {
			u.Log.Errorf("unauthorized access attempt: %v", err)
			return nil, errors.New(http.StatusText(http.StatusForbidden))
		}

		response := converter.TodoToResponse(todoData, isAdmin)

		if err := u.Cache.Set(key, response, 5*time.Minute); err != nil {
			u.Log.Errorf("failed to set data to cache: %v", err)
		}

		return response, nil
	}

	return data, nil
}

func (u *TodoUsecaseImpl) Search(ctx context.Context, request *model.TodoSearchRequest) (*model.Response[[]*model.TodoResponse], error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.helper.GetJWTClaims(ctx)
	if err != nil {
		u.Log.Errorf("failed to get JWT claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	isAdmin := u.helper.IsAdmin(ctx)
	userID := claims.UserID

	opts := model.TodoQueryOptions{
		Page:    request.Page,
		Size:    request.Size,
		IsAdmin: isAdmin,
		Title:   &request.Title,
	}

	// If not admin, always filter by user's ID
	if !isAdmin {
		opts.UserID = &userID
	}

	// Ensure valid pagination parameters
	if request.Size <= 0 {
		request.Size = 10 // Default page size
	}
	if request.Page <= 0 {
		request.Page = 1 // Default page number
	}

	// Try to get cached data
	cacheKey := u.helper.BuildCacheKey(opts)
	var cachedResponse *model.Response[[]*model.TodoResponse]
	if err := u.Cache.Get(cacheKey, &cachedResponse); err == nil {
		return cachedResponse, nil
	}

	// If cache miss, get from database
	var todos []entity.Todo
	totalItems, err := u.TodoRepository.GetPaginated(u.DB, &todos, opts)
	if err != nil {
		u.Log.Errorf("failed to get todos: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if totalItems == 0 || len(todos) == 0 {
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	response := converter.TodosToPaginatedResponse(todos, totalItems, request.Page, request.Size, isAdmin)

	// Cache the response
	if err := u.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		u.Log.Errorf("failed to cache response: %v", err)
	}

	return response, nil
}

func (u *TodoUsecaseImpl) GetAll(ctx context.Context, request *model.TodoGetAllRequest) (*model.Response[[]*model.TodoResponse], error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.helper.GetJWTClaims(ctx)
	if err != nil {
		u.Log.Errorf("failed to get JWT claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	isAdmin := u.helper.IsAdmin(ctx)
	userID := claims.UserID

	opts := model.TodoQueryOptions{
		Page:    request.Page,
		Size:    request.Size,
		IsAdmin: isAdmin,
	}

	// If not admin, always filter by user's ID
	if !isAdmin {
		opts.UserID = &userID
	}

	// Ensure valid pagination parameters
	if request.Size <= 0 {
		request.Size = 10 // Default page size
	}
	if request.Page <= 0 {
		request.Page = 1 // Default page number
	}

	// Try to get cached data
	cacheKey := u.helper.BuildCacheKey(opts)
	var cachedResponse *model.Response[[]*model.TodoResponse]
	if err := u.Cache.Get(cacheKey, &cachedResponse); err == nil {
		return cachedResponse, nil
	}

	// If cache miss, get from database
	var todos []entity.Todo
	totalItems, err := u.TodoRepository.GetPaginated(u.DB, &todos, opts)
	if err != nil {
		u.Log.Errorf("failed to get todos: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if totalItems == 0 || len(todos) == 0 {
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	response := converter.TodosToPaginatedResponse(todos, totalItems, request.Page, request.Size, isAdmin)

	// Cache the response
	if err := u.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		u.Log.Errorf("failed to cache response: %v", err)
	}

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
