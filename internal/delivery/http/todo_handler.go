package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models"
	"github.com/savioruz/mikti-task/tree/week-4/internal/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TodoHandler struct {
	Log     *logrus.Logger
	UseCase *usecases.TodoUsecase
}

func NewTodoHandler(log *logrus.Logger, useCase *usecases.TodoUsecase) *TodoHandler {
	return &TodoHandler{
		Log:     log,
		UseCase: useCase,
	}
}

// Create function is a handler to create a new todo
// @Summary Create a new todo
// @Description Create a new todo
// @Tags todo
// @Accept json
// @Produce json
// @Param todo body models.TodoCreateRequest true "Todo data"
// @Success 201 {object} models.ResponseSuccess[models.TodoResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @security ApiKeyAuth
// @Router /todo [post]
func (h *TodoHandler) Create(ctx echo.Context) error {
	request := new(models.TodoCreateRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create todo: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusCreated, models.ResponseSuccess[*models.TodoResponse]{
		Data: response,
	})
}

// GetByID function is a handler to get todo by ID
// @Summary Get todo by ID
// @Description Get todo by ID
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.ResponseSuccess[models.TodoResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @security ApiKeyAuth
// @Router /todo/{id} [get]
func (h *TodoHandler) GetByID(ctx echo.Context) error {
	request := new(models.TodoGetRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.Get(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get todo: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, models.ResponseSuccess[*models.TodoResponse]{
		Data: response,
	})
}

// Update function is a handler to update todo
// @Summary Update todo
// @Description Update todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body models.TodoUpdateRequest true "Todo data"
// @Success 200 {object} models.ResponseSuccess[models.TodoResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @security ApiKeyAuth
// @Router /todo/{id} [put]
func (h *TodoHandler) Update(ctx echo.Context) error {
	id := &models.TodoUpdateIDRequest{
		ID: ctx.Param("id"),
	}

	request := new(models.TodoUpdateRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request body: %v", err)
		return HandleError(ctx, http.StatusBadRequest, ErrorBindingRequest)
	}

	response, err := h.UseCase.Update(ctx.Request().Context(), id, request)
	if err != nil {
		h.Log.Errorf("failed to update todo: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, models.ResponseSuccess[*models.TodoResponse]{
		Data: response,
	})
}

// Delete function is a handler to delete todo
// @Summary Delete todo
// @Description Delete todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.ResponseSuccess[models.TodoResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @security ApiKeyAuth
// @Router /todo/{id} [delete]
func (h *TodoHandler) Delete(ctx echo.Context) error {
	request := new(models.TodoDeleteRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	err := h.UseCase.Delete(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to delete todo: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, models.ResponseSuccess[*models.TodoResponse]{
		Data: nil,
		Message: &models.Message{
			Message: fmt.Sprintf("todo with ID %s has been deleted", request.ID),
		},
	})
}

// List function is a handler to list todo
// @Summary List todo
// @Description List todo
// @Tags todo
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} models.ResponseSuccess[[]models.TodoResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @security ApiKeyAuth
// @Router /todo [get]
func (h *TodoHandler) List(ctx echo.Context) error {
	request := new(models.TodoListRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.List(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to list todo: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, models.ResponseSuccess[[]*models.TodoResponse]{
		Data: response,
	})
}
