package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/models"
)

var (
	ErrorBindingRequest = errors.New("failed to bind request")
	ErrorValidation     = errors.New("validation error")
	ErrorInternalServer = errors.New("failed to process request")
	ErrorNotFound       = errors.New("data not found")
	ErrorConflict       = errors.New("data already exists")
)

func HandleError(c echo.Context, status int, err error) error {
	return c.JSON(status, models.WebResponse[any]{
		Error: err.Error(),
	})
}
