package bootstrap

import (
	"strings"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() (*Env, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !strings.Contains(err.Error(), "no such file") {
			return nil, err
		}
	}

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}

func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("CONTEXT_TIMEOUT", 30)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "pixel")
	viper.SetDefault("DB_PASS", "123456")
	viper.SetDefault("DB_NAME", "omni_pixel")
	viper.SetDefault("ACCESS_TOKEN_EXPIRY_HOUR", 2)
	viper.SetDefault("REFRESH_TOKEN_EXPIRY_HOUR", 168)
	viper.SetDefault("ACCESS_TOKEN_SECRET", "change-me-in-production")
}
