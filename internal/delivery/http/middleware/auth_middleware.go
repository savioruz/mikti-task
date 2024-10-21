package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/auth"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
			}

			claims, err := jwtService.ValidateToken(bearerToken[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
