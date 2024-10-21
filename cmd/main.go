package main

import (
	"fmt"
	_ "github.com/savioruz/mikti-task/tree/post-1/docs"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/config"
	"time"
)

// @title			Bookstore API
// @version		1.0
// @description	This is a simple bookstore API
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	jakueenak@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes		http
// @basePath		/api/v1
func main() {
	log := config.NewLogrus()
	viper := config.NewViper()
	db := config.NewDatabase(log, "./tmp/books.db")
	validate := config.NewValidator()
	app, log := config.NewEcho()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})

	port := viper.GetString("APP_PORT")
	go func() {
		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatal("shutting down the server")
		}
	}()

	config.GracefulShutdown(app, log, 10*time.Second)
}
