package config

import (
	"github.com/spf13/viper"
)

// GetEnv - Get environment variable
func GetEnv(environmentVariableName string, configFilePath string) *string {
	// First, try to get environment variable. It's first.
	var value, ok = viper.Get(environmentVariableName).(string)
	if ok {
		return &value
	}

	// Second, try to find value from config file. Must check IsSet
	if viper.IsSet(configFilePath) {
		value = viper.GetString(configFilePath)
		return &value
	}

	return nil
}
