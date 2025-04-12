package config

import "os"

type ConfigApp struct {
	Port          string `mapstructure:"port" json:"port" yaml:"port"`
	DbConnect     string `mapstructure:"db_connect" json:"db_connect" yaml:"db_connect"`
	DbConnectRead string `mapstructure:"db_connect_read" json:"db_connect_read" yaml:"db_connect_read"`
	RedisConnect  string `mapstructure:"redis_connect" json:"redis_connect" yaml:"redis_connect"`
	RedisPass     string `mapstructure:"redis_pass" json:"redis_pass" yaml:"redis_pass"`
	DbName        string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	SecretKey     string `mapstructure:"private_key" json:"private_key" yaml:"private_key"`
	PoliciesPath  *map[string][]string
}

func NewConfigEnv() *ConfigApp {
	return &ConfigApp{
		Port:          os.Getenv("SERVER_PORT"),
		RedisConnect:  os.Getenv("REDIS_CONNECTION"),
		RedisPass:     os.Getenv("REDIS_PASS"),
		DbConnect:     os.Getenv("DB_CONNECTION"),
		DbConnectRead: os.Getenv("DB_CONNECTION_READ"),
		DbName:        os.Getenv("DB_NAME"),
		SecretKey:     os.Getenv("SECRET_KEY"),
	}
}
