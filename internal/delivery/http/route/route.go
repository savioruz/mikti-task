package route

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/savioruz/mikti-task/internal/delivery/graph/handler"
	"github.com/savioruz/mikti-task/internal/delivery/http/handler/todo"
	"github.com/savioruz/mikti-task/internal/delivery/http/handler/user"
	swagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	App            *echo.Echo
	GraphQLHandler *handler.GraphQLHandler
	TodoHandler    *todo.TodoHandlerImpl
	UserHandler    *user.UserHandlerImpl
	AuthMiddleware echo.MiddlewareFunc
}

func (c *Config) Setup() {
	c.publicRoutes()
	c.protectedRoutes()
	c.graphqlRoutes()
	c.swaggerRoutes()
	c.App.Use(middleware.Recover())
}

func (c *Config) publicRoutes() {
	g := c.App.Group("/api/v1")
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	g.POST("/users", c.UserHandler.Register)
	g.POST("/users/login", c.UserHandler.Login)
	g.POST("/users/refresh", c.UserHandler.Refresh)
}

func (c *Config) protectedRoutes() {
	g := c.App.Group("/api/v1")
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	g.Use(c.AuthMiddleware)
	g.POST("/todo", c.TodoHandler.Create)
	g.GET("/todo", c.TodoHandler.GetAll)
	g.GET("/todo/search", c.TodoHandler.Search)
	g.GET("/todo/:id", c.TodoHandler.GetByID)
	g.PUT("/todo/:id", c.TodoHandler.Update)
	g.DELETE("/todo/:id", c.TodoHandler.Delete)
}

func (c *Config) graphqlRoutes() {
	g := c.App.Group("/api/v1/graphql")
	g.Use(c.AuthMiddleware)
	g.POST("", c.GraphQLHandler.GraphQLHandler)
	c.App.GET("/playground", c.GraphQLHandler.PlaygroundHandler)
}

func (c *Config) swaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}
