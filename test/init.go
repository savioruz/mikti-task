package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/week-4/config"
	"github.com/savioruz/mikti-task/tree/week-4/internal/cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	app      *echo.Echo
	db       *gorm.DB
	redis    *cache.Cache
	log      *logrus.Logger
	validate *validator.Validate
	c        *viper.Viper
)

func init() {
	c = config.NewViper()
	log = config.NewLogrus()
	validate = config.NewValidator()
	db = config.NewDatabase(c, log)
	redis = config.NewRedisClient(c, log)
	jwt := config.NewJWT(c)

	var newLog *logrus.Logger
	app, newLog = config.NewEcho()
	if newLog != nil {
		log = newLog
	}

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Cache:    redis,
		App:      app,
		Log:      log,
		Validate: validate,
		JWT:      jwt,
	})
}
