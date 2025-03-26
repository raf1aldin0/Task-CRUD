package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	RedisHost  string
	RedisPort  string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		DbHost:     viper.GetString("DB_HOST"),
		DbPort:     viper.GetString("DB_PORT"),
		DbUser:     viper.GetString("DB_USER"),
		DbPassword: viper.GetString("DB_PASSWORD"),
		DbName:     viper.GetString("DB_NAME"),
		RedisHost:  viper.GetString("REDIS_HOST"),
		RedisPort:  viper.GetString("REDIS_PORT"),
	}
}
