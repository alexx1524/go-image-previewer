package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServer
	Storage    Storage
	Cache      Cache
	Log        Log
}

type HTTPServer struct {
	Address           string `mapstructure:"Address"`
	Port              int    `mapstructure:"Port"`
	ReadHeaderTimeout int    `mapstructure:"ReadHeaderTimeout"`
}

type Storage struct {
	ImagesPath string `mapstructure:"ImagesPath"`
}

type Cache struct {
	Mode     string `mapstructure:"Mode"`
	LRUCache LRUCache
}

type LRUCache struct {
	ItemsCount int `mapstructure:"ItemsCount"`
}

type Log struct {
	LogFile string `mapstructure:"LogFile"`
	Level   string `mapstructure:"Level"`
}

func NewConfig(configFilePath string) (Config, error) {
	return buildConfig(configFilePath)
}

func buildConfig(configFilePath string) (Config, error) {
	var config Config

	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
