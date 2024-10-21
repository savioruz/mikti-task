package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/controllers"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/models"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/routes"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/utils"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/utils/usecases"
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
	var query *utils.Query[models.Book]

	// Setup UseCase
	bookUseCase := usecases.NewBookUseCase(config.DB, config.Log, config.Validate, query)

	// Setup Controller
	bookController := controllers.NewBookController(config.Log, bookUseCase)

	// Setup Routes
	routeConfig := routes.Config{
		App:            config.App,
		BookController: bookController,
	}

	routeConfig.Setup()

	config.Log.Infof("Application is ready to serve")
}
