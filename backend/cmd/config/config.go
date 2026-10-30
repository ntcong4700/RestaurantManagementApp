package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{
		DBHost:     viper.GetString("DBHost"),
		DBPort:     viper.GetInt("DBPort"),
		DBUser:     viper.GetString("DBUser"),
		DBPassword: viper.GetString("DBPassword"),
		DBName:     viper.GetString("DBName"),
	}
	return config, nil
}
