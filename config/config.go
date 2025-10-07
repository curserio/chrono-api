package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AppConfig struct {
	Name            string
	DevMode         bool
	DefaultLanguage string
}

func NewConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Config file not found: %v, using env only", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Addr: v.GetString("server.addr"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			DBName:   v.GetString("database.dbname"),
			SSLMode:  v.GetString("database.sslmode"),
		},
		App: AppConfig{
			Name:            v.GetString("app.name"),
			DevMode:         v.GetBool("app.dev_mode"),
			DefaultLanguage: v.GetString("app.default_language"),
		},
	}
	return cfg
}
