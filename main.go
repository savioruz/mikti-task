package main

import (
	_ "github.com/savioruz/mikti-task/week-3/docs"
	"github.com/savioruz/mikti-task/week-3/internal/config"
)

// @title Todo API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func main() {
	log := config.NewLogrus()
	db := config.NewDatabase(log)
	validate := config.NewValidator()
	app := config.NewEcho()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})
}
