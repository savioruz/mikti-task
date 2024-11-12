package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/graph/handler"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/graph/resolvers"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/handler/todo"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/handler/user"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/middleware"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/route"
	"github.com/savioruz/mikti-task/tree/week-4/internal/platform/cache"
	"github.com/savioruz/mikti-task/tree/week-4/internal/platform/jwt"
	todoRepo "github.com/savioruz/mikti-task/tree/week-4/internal/repositories/todo"
	userRepo "github.com/savioruz/mikti-task/tree/week-4/internal/repositories/user"
	todoUsecase "github.com/savioruz/mikti-task/tree/week-4/internal/usecases/todo"
	userUsecase "github.com/savioruz/mikti-task/tree/week-4/internal/usecases/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Cache    *cache.ImplCache
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
	JWT      *jwt.JWTConfig
}

func Bootstrap(config *BootstrapConfig) error {
	// Initialize repositories
	todoRepository := todoRepo.NewTodoRepository(config.DB, config.Log)
	userRepository := userRepo.NewUserRepository(config.DB, config.Log)

	// Initialize JWT service
	jwtService := jwt.NewJWTService(config.JWT)

	// Initialize usecases
	todoUC := todoUsecase.NewTodoUsecaseImpl(
		config.DB,
		config.Cache,
		config.Log,
		config.Validate,
		todoRepository,
	)

	userUC := userUsecase.NewUserUsecaseImpl(
		config.DB,
		config.Log,
		config.Validate,
		userRepository,
		jwtService,
	)

	// Initialize handlers
	todoHandler := todo.NewTodoHandlerImpl(config.Log, todoUC)
	userHandler := user.NewUserHandlerImpl(config.Log, userUC)

	// Initialize GraphQL
	resolver := resolvers.NewResolver(todoUC)
	graphQLHandler := handler.NewGraphQLHandler(resolver)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Setup routes
	routeConfig := &route.Config{
		App:            config.App,
		GraphQLHandler: graphQLHandler,
		TodoHandler:    todoHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()

	config.Log.Info("Application is ready")
	return nil
}
