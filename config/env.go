package config

import (
	"github.com/spf13/viper"
)

// GetEnv - Get environment variable
func GetEnv(name string) *string {
	value, ok := viper.Get(name).(string)
	if ok {
		return &value
	}

	return nil
}
