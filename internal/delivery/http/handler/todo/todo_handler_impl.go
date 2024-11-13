package todo

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/internal/delivery/http/handler"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/savioruz/mikti-task/internal/usecases/todo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TodoHandlerImpl struct {
	Log  *logrus.Logger
	Todo todo.TodoUsecase
}

func NewTodoHandlerImpl(log *logrus.Logger, t todo.TodoUsecase) *TodoHandlerImpl {
	return &TodoHandlerImpl{
		Log:  log,
		Todo: t,
	}
}

// Create function is a handler to create a new todo
// @Summary Create a new todo
// @Description Create a new todo
// @Tags todo
// @Accept json
// @Produce json
// @Param todo body model.TodoCreateRequest true "Todo data"
// @Success 201 {object} model.Response[model.TodoResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @security ApiKeyAuth
// @Router /todo [post]
func (h *TodoHandlerImpl) Create(ctx echo.Context) error {
	request := new(model.TodoCreateRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.Todo.Create(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create todo: %v", err)
		if err.Error() == "Bad Request" {
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		} else {
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// GetByID function is a handler to get todo by ID
// @Summary Get todo by ID
// @Description Get todo by ID
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} model.Response[model.TodoResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @security ApiKeyAuth
// @Router /todo/{id} [get]
func (h *TodoHandlerImpl) GetByID(ctx echo.Context) error {
	request := new(model.TodoGetRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.Todo.Get(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get todo: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Not Found":
			return handler.HandleError(ctx, http.StatusNotFound, handler.ErrNotFound)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// Update function is a handler to update todo
// @Summary Update todo
// @Description Update todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body model.TodoUpdateRequest true "Todo data"
// @Success 200 {object} model.Response[model.TodoResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @security ApiKeyAuth
// @Router /todo/{id} [put]
func (h *TodoHandlerImpl) Update(ctx echo.Context) error {
	id := &model.TodoUpdateIDRequest{
		ID: ctx.Param("id"),
	}

	request := new(model.TodoUpdateRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request body: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.Todo.Update(ctx.Request().Context(), id, request)
	if err != nil {
		h.Log.Errorf("failed to update todo: %v", err)
		if err.Error() == "Bad Request" {
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		} else {
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// Delete function is a handler to delete todo
// @Summary Delete todo
// @Description Delete todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} model.Response[model.TodoResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @security ApiKeyAuth
// @Router /todo/{id} [delete]
func (h *TodoHandlerImpl) Delete(ctx echo.Context) error {
	request := new(model.TodoDeleteRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	_, err := h.Todo.Delete(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to delete todo: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Not Found":
			return handler.HandleError(ctx, http.StatusNotFound, handler.ErrNotFound)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetAll function is a handler to list todo
// @Summary List todo
// @Description List todo
// @Tags todo
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} model.Response[[]model.TodoResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @security ApiKeyAuth
// @Router /todo [get]
func (h *TodoHandlerImpl) GetAll(ctx echo.Context) error {
	request := new(model.TodoGetAllRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.Todo.GetAll(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to list todo: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Not Found":
			return handler.HandleError(ctx, http.StatusNotFound, handler.ErrNotFound)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
