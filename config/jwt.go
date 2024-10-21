package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type JWTConfig struct {
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

func NewJWT(viper *viper.Viper) *JWTConfig {
	fmt.Println("JWT_SECRET", viper.GetString("JWT_SECRET"))
	fmt.Println("JWT_ACCESS_EXPIRY", viper.GetDuration("JWT_ACCESS_EXPIRY"))
	fmt.Println("JWT_REFRESH_EXPIRY", viper.GetDuration("JWT_REFRESH_EXPIRY"))
	return &JWTConfig{
		JWTSecret:        viper.GetString("JWT_SECRET"),
		JWTAccessExpiry:  viper.GetDuration("JWT_ACCESS_EXPIRY"),
		JWTRefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),
	}
}
