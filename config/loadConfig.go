package config

import "github.com/spf13/viper"

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
