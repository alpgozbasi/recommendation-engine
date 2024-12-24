package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// appconfig will hold  all configuration values after loading from config.yaml

type AppConfig struct {
	App      AppSettings
	Database DatabaseSettings
	Redis    RedisSettings
}

type AppSettings struct {
	Port int
}

type DatabaseSettings struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisSettings struct {
	Host     string
	Password string
	DB       int
}

// loadconfig loads configuration from file and env. variables
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // in root dir

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error reading config file, %s", err)
		return nil, err
	}

	var cfg AppConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &cfg, nil
}
