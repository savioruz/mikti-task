package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-3/internal/delivery/restful"
	"github.com/savioruz/mikti-task/tree/week-3/internal/delivery/restful/middleware"
	"github.com/savioruz/mikti-task/tree/week-3/internal/delivery/restful/route"
	"github.com/savioruz/mikti-task/tree/week-3/internal/repositories"
	"github.com/savioruz/mikti-task/tree/week-3/internal/usecases"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	todoRepository := repositories.NewTodoRepository(config.Log)
	userRepository := repositories.NewUserRepository(config.Log)

	// Setup usecases
	todoUsecase := usecases.NewTodoUsecase(config.DB, config.Log, config.Validate, todoRepository)
	userUsecase := usecases.NewUserUsecase(config.DB, config.Log, config.Validate, userRepository)

	// Setup handlers
	todoHandler := restful.NewTodoHandler(config.Log, todoUsecase)
	userHandler := restful.NewUserHandler(config.Log, userUsecase)

	// Setup middleware
	authMiddleware := middleware.NewAuth(userUsecase)

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
