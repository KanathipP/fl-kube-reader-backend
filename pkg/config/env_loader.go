// Package config use for getting configurations
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GoDotEnvConfig struct {
	KubeConfig   KubeConfig
	ServerConfig ServerConfig
}

type KubeConfig struct {
	Namespace     string
	LabelSelector string
}

type ServerConfig struct {
	FrontendURL string
	Port        string
}

func LoadDotEnvConfig() GoDotEnvConfig {
	_ = godotenv.Load()

	kubeConfig := KubeConfig{
		Namespace:     Get("NAMESPACE", "flwr"),
		LabelSelector: Get("LABEL_SELECTOR", "app.kubernetes.io/component in (clientapp,serverapp)"),
	}

	serverConfig := ServerConfig{
		FrontendURL: Get("FRONTEND_URL", "http://localhost:8080"),
		Port:        Get("PORT", "3000"),
	}

	config := GoDotEnvConfig{
		KubeConfig:   kubeConfig,
		ServerConfig: serverConfig,
	}

	return config
}

func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err == nil {
			return n
		}
	}
	return defaultValue
}
