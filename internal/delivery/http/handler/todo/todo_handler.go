package todo

import (
	"github.com/labstack/echo/v4"
)

type TodoHandler interface {
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	GetByID(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	Delete(ctx echo.Context) error
}
