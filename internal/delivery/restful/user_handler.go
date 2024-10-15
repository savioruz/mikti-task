package restful

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
	"github.com/savioruz/mikti-task/tree/week-3/internal/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	Log     *logrus.Logger
	UseCase *usecases.UserUsecase
}

func NewUserHandler(log *logrus.Logger, useCase *usecases.UserUsecase) *UserHandler {
	return &UserHandler{
		Log:     log,
		UseCase: useCase,
	}
}

// Register function is a handler to register a new user
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param register body models.RegisterUserRequest true "Register User Request"
// @Success 201 {object} models.ResponseSuccess[models.UserResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /users [post]
func (h *UserHandler) Register(ctx echo.Context) error {
	request := new(models.RegisterUserRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to register user: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusCreated, models.ResponseSuccess[*models.UserResponse]{
		Data: response,
	})
}

// Login function is a handler to login a user
// @Summary Login a user
// @Description Login a user
// @Tags User
// @Accept json
// @Produce json
// @Param login body models.LoginUserRequest true "Login User Request"
// @Success 200 {object} models.ResponseSuccess[models.UserResponse]
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /users/login [post]
func (h *UserHandler) Login(ctx echo.Context) error {
	request := new(models.LoginUserRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.Login(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to login user: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, models.ResponseSuccess[*models.UserResponse]{
		Data: response,
	})
}

func (h *UserHandler) Logout(ctx echo.Context) error {
	request := new(models.LogoutUserRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return HandleError(ctx, http.StatusInternalServerError, ErrorBindingRequest)
	}

	response, err := h.UseCase.Logout(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to logout user: %v", err)
		if err.Error() == "Bad Request" {
			return HandleError(ctx, http.StatusBadRequest, ErrValidation)
		} else {
			return HandleError(ctx, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusNoContent, models.ResponseSuccess[bool]{Data: response})
}
