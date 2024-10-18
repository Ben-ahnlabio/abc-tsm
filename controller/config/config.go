package config

import (
	"os"
)

type Config struct {
	PlayerIndex          string `env:"PLAYER_INDEX"`
	AppName              string `env:"APP_NAME"`
	AppVersion           string `env:"APP_VERSION"`
	BuildType            string `env:"BUILD_TYPE"`
	NodeUrl              string `env:"NODE_URL"`
	NodeApiKey           string `env:"NODE_API_KEY"`
	NodePubicKey         string `env:"NODE_PUBLIC_KEY"`
	AnotherNodePublicKey string `env:"ANOTHER_NODE_PUBLIC_KEY"`
}

func GetConfig() *Config {
	return &Config{
		AppName:              os.Getenv("APP_NAME"),
		AppVersion:           os.Getenv("APP_VERSION"),
		BuildType:            os.Getenv("BUILD_TYPE"),
		PlayerIndex:          os.Getenv("PLAYER_INDEX"),
		NodeUrl:              os.Getenv("NODE_URL"),
		NodeApiKey:           os.Getenv("NODE_API_KEY"),
		NodePubicKey:         os.Getenv("NODE_PUBLIC_KEY"),
		AnotherNodePublicKey: os.Getenv("ANOTHER_NODE_PUBLIC_KEY"),
	}
}
