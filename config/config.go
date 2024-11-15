package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DPI         int     `json:"dpi"`
	Sensitivity float64 `json:"sensitivity"`
}

func defaultConfig() Config {
	return Config{
		DPI:         400,
		Sensitivity: 0.08,
	}
}

func loadConfig(filename string) *Config {
	defaultConfig := defaultConfig()

	data, err := os.ReadFile(filename)
	if err != nil {
		return &defaultConfig
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return &defaultConfig
	}

	return &config
}

var runtimeConfig *Config = nil

func GetConfig() Config {
	if runtimeConfig == nil {
		runtimeConfig = loadConfig("config.json")
	}

	return *runtimeConfig

}
