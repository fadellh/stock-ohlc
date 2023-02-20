package config

import (
	"os"

	"github.com/spf13/viper"
)

type ContextKey string

type Config struct {
	AppEnv          string `mapstructure:"APP_ENV"`
	AppIsDev        bool
	RedisConnection string `mapstructure:"REDIS_CONNECTION"`
	RedisAddress    string `mapstructure:"REDIS_ADDRESS"`
	RedisUsername   string `mapstructure:"REDIS_USERNAME"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase   int    `mapstructure:"REDIS_DATABASE"`
	RedisExpired    int    `mapstructure:"REDIS_EXPIRED"`
}

func NewConfig() (*Config, error) {
	env := os.Getenv("APPENV")
	if env == "" {
		env = "local"
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("summary-service/config")
	viper.SetConfigName(env)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigName("placeholder")

			if err := viper.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.AppIsDev = cfg.AppEnv == "staging" || cfg.AppEnv == "local" || cfg.AppEnv == "dev"

	return cfg, nil
}
