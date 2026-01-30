package main

import (
	"log/slog"

	"todo-n8n/config"
)

func main() {
	apiUrl, err := config.GetEnv("API_URL")
	header, err := config.GetEnv("AUTH_HEADER")
	key, err := config.GetEnv("AUTH_KEY")
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(apiUrl)
	slog.Info(header)
	slog.Info(key)
}
