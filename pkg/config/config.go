package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Read(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "./config.yaml"
	}
	viper.SetConfigFile(configPath)
	// Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Zusammensetzen der URL aus HOST und PORT
	config.RabbitMQ.URL = fmt.Sprintf("%s:%s", config.RabbitMQ.Host, config.RabbitMQ.Port)
	config.InfluxDB.URL = fmt.Sprintf("http://%s:%s", config.InfluxDB.Host, config.InfluxDB.Port)
	config.Redis.URL = fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	return &config, nil
}
