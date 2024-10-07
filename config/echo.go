package config

import (
	"github.com/labstack/echo/v4"
)

// NewEcho is a function to create new echo instance
func NewEcho() *echo.Echo {
	return echo.New()
}
