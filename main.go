package main

import (
	"fmt"
	"github.com/savioruz/mikti-task/tree/week-3/config"
	_ "github.com/savioruz/mikti-task/tree/week-3/docs"
	"time"
)

// @title Todo API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	viper := config.NewViper()
	log := config.NewLogrus()
	db := config.NewDatabase(viper, log)
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
