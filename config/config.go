package config

import (
	"fmt"
	"os"

	Viper "github.com/spf13/viper"
)

type Config struct {
	Url   string
	Mongo Mongos
}

var config *Config

func init() {
	var folder string

	env := os.Getenv("APPLICATION_ENV")

	switch env {
	case "local":
		folder = env
	default:
		folder = "local"
	}

	path := fmt.Sprintf("config/%v", folder)

	config = new(Config)
	fetchDataToConfig(path, "crawler", &config)

	fetchDataToConfig(path, "mongo", &config.Mongo)
}

func fetchDataToConfig(configPath, configName string, result interface{}) {
	viper := Viper.New()
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	err := viper.ReadInConfig() // Find and read the config file
	if err == nil {             // Handle errors reading the config file
		err = viper.Unmarshal(result)
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}

func GetConfig() *Config {
	return config
}
