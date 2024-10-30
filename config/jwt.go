package config

import (
	"github.com/spf13/viper"
	"time"
)

type JWTConfig struct {
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

func NewJWT(viper *viper.Viper) *JWTConfig {
	return &JWTConfig{
		JWTSecret:        viper.GetString("JWT_SECRET"),
		JWTAccessExpiry:  viper.GetDuration("JWT_ACCESS_EXPIRY"),
		JWTRefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),
	}
}
