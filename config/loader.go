package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadYAMLConfig loads the YAML config file
func LoadYAMLConfig(configFileName string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(configFileName)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file: %v", err)
	}

	return v, nil
}
