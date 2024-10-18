package config

import "os"

type Config struct {
	AppName    string `env:"APP_NAME"`
	AppVersion string `env:"APP_VERSION"`
	BuildType  string `env:"BUILD_TYPE"`
	Player1Url string `env:"PLAYER1_URL"`
	Player2Url string `env:"PLAYER2_URL"`
}

func GetConfig() *Config {
	return &Config{
		AppName:    os.Getenv("APP_NAME"),
		AppVersion: os.Getenv("APP_VERSION"),
		BuildType:  os.Getenv("BUILD_TYPE"),
		Player1Url: os.Getenv("PLAYER1_URL"),
		Player2Url: os.Getenv("PLAYER2_URL"),
	}
}
