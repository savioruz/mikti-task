package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/controllers"
	swagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	App            *echo.Echo
	BookController *controllers.BookController
}

func (c *Config) Setup() {
	c.PublicRoutes()
	c.swaggerRoutes()
	c.App.Use(middleware.Recover())
}

func (c *Config) PublicRoutes() {
	g := c.App.Group("/api/v1")
	g.POST("/book", c.BookController.CreateBook)
	g.GET("/book/:id", c.BookController.GetBookByID)
	g.GET("/books", c.BookController.GetAllBooks)
	g.PUT("/book/:id", c.BookController.UpdateBook)
	g.DELETE("/book/:id", c.BookController.DeleteBook)
}

func (c *Config) swaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}
