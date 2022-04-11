package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	InterfaceName string    `mapstructure:"interface"`
	Services      []Service `mapstructure:"services"`
	Teams         []Team    `mapstructure:"teams"`
}

type Service struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type Team struct {
	Name string `mapstructure:"name"`
	IP   string `mapstructure:"ip"`
}

func GetConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	var c Config

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &c, fmt.Errorf("%s", "config doesn't exist")
		}
		return &c, fmt.Errorf("%s", "read config error")
	}

	if err := viper.Unmarshal(&c); err != nil {
		return &c, fmt.Errorf("%s", "Config error: unmarshal error")
	}

	if c.InterfaceName == "" {
		return &c, fmt.Errorf("%s", "Config error: set interface")
	}

	if len(c.Services) == 0 {
		return &c, fmt.Errorf("%s", "Config error: set services")
	} else {
		for _, s := range c.Services {
			if s.Name == "" {
				return &c, fmt.Errorf("%s", "Config error: set name")
			}
		}
	}

	if len(c.Teams) == 0 {
		return &c, fmt.Errorf("%s", "Config error: set teams")
	}

	return &c, nil
}
