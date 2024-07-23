package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ALLOWED_ORIGINS string `mapstructure:"ALLOWED_ORIGINS"`
	PORT            string `mapstructure:"PORT"`

	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBUsername string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBTestName string `mapstructure:"POSTGRES_DB_TEST"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`

	MONGODB_HOST   string `mapstructure:"MONGODB_HOST"`
	MONGODB_DB     string `mapstructure:"MONGODB_DB"`
	MONGODB_PORT   string `mapstructure:"MONGODB_PORT"`
	MONGODB_STRING string `mapstructure:"MONGODB_STRING"`

	TokenSecret    string        `mapstructure:"TOKEN_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
