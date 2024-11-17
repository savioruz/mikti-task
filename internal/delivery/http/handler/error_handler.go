package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

var (
	ErrorBindingRequest = errors.New("failed to bind request")
	ErrValidation       = errors.New("validation error")
	ErrorInternalServer = errors.New("failed to process request")
	ErrorUnauthorized   = errors.New("unauthorized")
	ErrorConflict       = errors.New("conflict")
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
)

func HandleError(c echo.Context, status int, err error) error {
	return c.JSON(status, model.NewErrorResponse[any](status, err.Error()))
}
