package route

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/savioruz/mikti-task/tree/week-3/internal/delivery/restful"
	swagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	App            *echo.Echo
	TodoHandler    *restful.TodoHandler
	UserHandler    *restful.UserHandler
	AuthMiddleware echo.MiddlewareFunc
}

func (c *Config) Setup() {
	c.publicRoutes()
	c.guestRoutes()
	c.swaggerRoutes()
	c.App.Use(middleware.Recover())
}

func (c *Config) guestRoutes() {
	g := c.App.Group("/api/v1")
	g.POST("/users", c.UserHandler.Register)
	g.POST("/users/login", c.UserHandler.Login)
}

func (c *Config) publicRoutes() {
	g := c.App.Group("/api/v1")
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	g.Use(c.AuthMiddleware)
	g.POST("/todo", c.TodoHandler.Create)
	g.GET("/todo", c.TodoHandler.List)
	g.GET("/todo/:id", c.TodoHandler.GetByID)
	g.PUT("/todo/:id", c.TodoHandler.Update)
	g.DELETE("/todo/:id", c.TodoHandler.Delete)
}

func (c *Config) swaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}
