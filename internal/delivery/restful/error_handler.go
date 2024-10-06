package restful

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/week-3/internal/models"
)

var (
	ErrorBindingRequest = errors.New("failed to bind request")
	ErrValidation       = errors.New("validation error")
	ErrorInternalServer = errors.New("failed to process request")
)

func HandleError(c echo.Context, status int, err error) error {
	return c.JSON(status, models.ResponseError{
		Error: err.Error(),
	})
}
