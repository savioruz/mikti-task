package config

import "github.com/spf13/viper"

func NewViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}
