package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/savioruz/mikti-task/internal/platform/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtService jwt.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errMessage := func(message string) error {
				return echo.NewHTTPError(http.StatusUnauthorized, model.NewErrorResponse[any](http.StatusUnauthorized, message))
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errMessage("Missing authorization header")
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				return errMessage("Invalid authorization header")
			}

			claims, err := jwtService.ValidateToken(bearerToken[1])
			if err != nil {
				return errMessage("Invalid token")
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
