package config

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Port is the HTTP port the server listens on.
	Port string
	// CORSAllowedOrigins lists origins allowed to call the API from a browser.
	// Empty disables CORS; "*" allows any origin.
	CORSAllowedOrigins []string
}

// LoadConfig reads configuration from environment variables, falling back to a
// .env file in the current directory if it exists. Environment variables win.
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("PORT", "8080")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read .env: %w", err)
		}
		// No .env file: rely on environment variables only.
	}

	config := &Config{
		DBHost:             v.GetString("DB_HOST"),
		DBPort:             v.GetInt("DB_PORT"),
		DBUser:             v.GetString("DB_USERNAME"),
		DBPassword:         v.GetString("DB_PASSWORD"),
		DBName:             v.GetString("DB_NAME"),
		Port:               v.GetString("PORT"),
		CORSAllowedOrigins: splitList(v.GetString("CORS_ALLOWED_ORIGINS")),
	}

	var missing []string
	if config.DBUser == "" {
		missing = append(missing, "DB_USERNAME")
	}
	if config.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	return config, nil
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
