package user

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/handler"
	"github.com/savioruz/mikti-task/tree/week-4/internal/domain/model"
	"github.com/savioruz/mikti-task/tree/week-4/internal/usecases/user"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandlerImpl struct {
	Log  *logrus.Logger
	User user.UserUsecase
}

func NewUserHandlerImpl(log *logrus.Logger, user user.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		Log:  log,
		User: user,
	}
}

// Register function is a handler to register a new user
// @Summary Register a new user
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User data"
// @Success 201 {object} model.Response[model.UserResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @Router /users [post]
func (h *UserHandlerImpl) Register(ctx echo.Context) error {
	request := new(model.RegisterRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.User.Create(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to register user: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Conflict":
			return handler.HandleError(ctx, http.StatusConflict, handler.ErrorConflict)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusCreated, model.Response[*model.UserResponse]{
		Data: response,
	})
}

// Login function is a handler to login a user
// @Summary Login a user
// @Description Login a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "User data"
// @Success 200 {object} model.Response[model.UserResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @Router /users/login [post]
func (h *UserHandlerImpl) Login(ctx echo.Context) error {
	request := new(model.LoginRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.User.Login(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to login user: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Unauthorized":
			return handler.HandleError(ctx, http.StatusUnauthorized, handler.ErrorUnauthorized)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, model.Response[*model.TokenResponse]{
		Data: response,
	})
}

// Refresh function is a handler to refresh token
// @Summary Refresh token
// @Description Refresh token
// @Tags user
// @Accept json
// @Produce json
// @Param token body model.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} model.Response[model.TokenResponse]
// @Failure 400 {object} model.Response[any]
// @Failure 500 {object} model.Response[any]
// @Router /users/refresh [post]
func (h *UserHandlerImpl) Refresh(ctx echo.Context) error {
	request := new(model.RefreshTokenRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrorBindingRequest)
	}

	response, err := h.User.RefreshToken(request)
	if err != nil {
		h.Log.Errorf("failed to refresh token: %v", err)
		switch {
		case err.Error() == "Bad Request":
			return handler.HandleError(ctx, http.StatusBadRequest, handler.ErrValidation)
		case err.Error() == "Unauthorized":
			return handler.HandleError(ctx, http.StatusUnauthorized, handler.ErrorUnauthorized)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, handler.ErrorInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, model.Response[*model.TokenResponse]{
		Data: response,
	})
}
