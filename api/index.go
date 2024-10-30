package handler

import (
	"github.com/savioruz/mikti-task/tree/week-4/config"
	_ "github.com/savioruz/mikti-task/tree/week-4/docs"
	"net/http"
)

// Handler is main function to run the application in vercel function
// @title Todo API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Handler(w http.ResponseWriter, r *http.Request) {
	viper := config.NewViper()
	log := config.NewLogrus()
	db := config.NewDatabase(viper, log)
	redis := config.NewRedisClient(viper, log)
	jwt := config.NewJWT(viper)
	validate := config.NewValidator()
	app, log := config.NewEcho()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Cache:    redis,
		App:      app,
		Log:      log,
		Validate: validate,
		JWT:      jwt,
	})

	app.ServeHTTP(w, r)
}
