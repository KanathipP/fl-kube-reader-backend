package main

import (
	"log/slog"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/config"
)

func main() {
	godotenv_env := config.Load()
	slog.Info("ENV has been loaded")
}
