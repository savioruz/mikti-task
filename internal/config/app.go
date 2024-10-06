package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/week-3/internal/delivery/restful"
	"github.com/savioruz/mikti-task/week-3/internal/delivery/restful/route"
	"github.com/savioruz/mikti-task/week-3/internal/repositories"
	"github.com/savioruz/mikti-task/week-3/internal/usecases"
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
	todoRepository := repositories.NewTodoRepository(config.DB, config.Log)

	// Setup usecases
	todoUsecase := usecases.NewTodoUsecase(config.DB, config.Log, config.Validate, todoRepository)

	// Setup handlers
	todoHandler := restful.NewTodoHandler(config.Log, todoUsecase)

	// Setup routes
	routeConfig := &route.Config{
		App:         config.App,
		TodoHandler: todoHandler,
	}
	routeConfig.Setup()

	err := config.App.Start(":3000")
	if err != nil {
		config.Log.Fatalf("failed to start application: %v", err)
		return
	}

	config.Log.Info("Application is ready")
}
