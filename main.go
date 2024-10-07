package main

import (
	"github.com/savioruz/mikti-task/tree/week-3/config"
	_ "github.com/savioruz/mikti-task/tree/week-3/docs"
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
	db := config.NewDatabase(log, "./tmp/todo.db")
	validate := config.NewValidator()
	app := config.NewEcho()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})

	app.Logger.Fatal(app.Start(":3000"))
}
