package user

import (
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
	Refresh(ctx echo.Context) error
}
