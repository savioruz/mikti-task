package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
	"github.com/savioruz/mikti-task/tree/week-3/internal/usecases"
	"net/http"
)

func NewAuth(userUseCase *usecases.UserUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				userUseCase.Log.Warnf("authorization token not found")
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			request := &models.VerifyUserRequest{Token: token}
			userUseCase.Log.Infof("verifying user with token: %s", token)

			auth, err := userUseCase.Verify(c.Request().Context(), request)
			if err != nil {
				userUseCase.Log.Errorf("failed to verify user: %v", err)
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			userUseCase.Log.Infof("user with token: %s is verified", token)
			c.Set("auth", auth)

			return next(c)
		}
	}
}

func GetUser(c echo.Context) *models.Auth {
	if auth, ok := c.Get("auth").(*models.Auth); ok {
		return auth
	}
	return nil
}
