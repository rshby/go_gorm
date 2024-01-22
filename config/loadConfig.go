package config

import "github.com/spf13/viper"

type Config struct {
	App      *App      `json:"app,omitempty"`
	Database *Database `json:"database,omitempty"`
}

type App struct {
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
	Port   int    `json:"port,omitempty"`
}

type Database struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
}

func LoadConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigType("json")
	config.SetConfigFile("config.json")
	config.AddConfigPath("./..")
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	// success get config
	return config, nil
}

func ConvertToObject(config *viper.Viper) *Config {
	return &Config{
		App: &App{
			Name:   config.GetString("app.name"),
			Author: config.GetString("app.author"),
			Port:   config.GetInt("app.port"),
		},
		Database: &Database{
			User:     config.GetString("database.user"),
			Password: config.GetString("database.password"),
			Host:     config.GetString("database.host"),
			Port:     config.GetInt("database.port"),
			Name:     config.GetString("database.name"),
		},
	}
}
