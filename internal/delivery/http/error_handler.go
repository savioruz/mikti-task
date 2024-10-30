package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models"
)

var (
	ErrorBindingRequest = errors.New("failed to bind request")
	ErrValidation       = errors.New("validation error")
	ErrorInternalServer = errors.New("failed to process request")
	ErrorUnauthorized   = errors.New("unauthorized")
)

func HandleError(c echo.Context, status int, err error) error {
	return c.JSON(status, models.ResponseError{
		Error: err.Error(),
	})
}
