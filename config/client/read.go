package config

import (
	"encoding/json"
	"os"
)

func ReadConfig(configPath string) (Config, error) {
	var cfg Config
	all, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	return cfg, json.Unmarshal(all, &cfg)
}

func MustReadConfig(configPath string) Config {
	cfg, err := ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return cfg
}
