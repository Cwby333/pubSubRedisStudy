package client

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	RedisConfig  RedisConfig  `yaml:"redis" yaml-required:"true"`
	ClientConfig ClientConfig `yaml:"client" yaml-required:"true"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" yaml-required:"true"`
	Username string `yaml:"username" yaml-required:"true"`
	Password string `yaml:"password" yaml-required:"true"`
	DB       int    `yaml:"db" yaml-required:"true"`
}

type ClientConfig struct {
	Username string `yaml:"username" yaml-required:"true"`
}

func MustLoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	configPath := os.Getenv("CHAT_CONFIG_PATH")
	if configPath == "" {
		panic("missing config path")
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
