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
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USERNAME"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
	}
	return config, nil
}
