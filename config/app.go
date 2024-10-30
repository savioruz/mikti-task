package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/cache"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/auth"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/middleware"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/route"
	"github.com/savioruz/mikti-task/tree/week-4/internal/repositories"
	"github.com/savioruz/mikti-task/tree/week-4/internal/usecases"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Cache    *cache.Cache
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
	JWT      *JWTConfig
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	todoRepository := repositories.NewTodoRepository(config.Log)
	userRepository := repositories.NewUserRepository(config.Log)

	// Setup jwt service
	jwtService := auth.NewJWTService(config.JWT.JWTSecret, config.JWT.JWTAccessExpiry, config.JWT.JWTRefreshExpiry)

	// Setup usecases
	todoUsecase := usecases.NewTodoUsecase(config.DB, config.Cache, config.Log, config.Validate, todoRepository)
	userUsecase := usecases.NewUserUsecase(config.DB, config.Log, config.Validate, userRepository, jwtService)

	// Setup handlers
	todoHandler := http.NewTodoHandler(config.Log, todoUsecase)
	userHandler := http.NewUserHandler(config.Log, userUsecase)

	// Setup middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Setup routes
	routeConfig := &route.Config{
		App:            config.App,
		TodoHandler:    todoHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()

	config.Log.Info("Application is ready")
}
