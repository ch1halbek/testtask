package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type DBConfig struct {
	DBPath string `json:"DBPath" yaml:"DBPath" mapstructure:"DBPath"`
}

func LoadDBConfig() (*DBConfig, error) {
	configFilePath := filepath.Join("internal", "config", "config.yaml")
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config DBConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config into struct: %v", err)
	}

	return &config, nil
}
