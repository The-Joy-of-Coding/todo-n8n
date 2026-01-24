package config

import (
	_ "embed"
	"log/slog"
	"os"
	"strings"
)

type ConfigDefault struct {
	URL string `json:"url"`
}

var Default ConfigDefault

func init() {
	setDefaults()
	setLogger()
}

func ApiUrl(logger *slog.Logger, url_key string) string {
	url := os.Getenv(url_key)
	if url == "" {
		logger.Error("No valid url found!")
		return Default.URL
	}
	if !strings.HasPrefix(url, "https://") {
		logger.Warn("Using insecure api endpoint.")
	}
	if !ping(url) {
		logger.Error("Url is invalid, please change the url to run the application.")
		return Default.URL
	}
	Default.URL = url
	return url
}

func ping(url string) bool {
	return url != ""
}

func setDefaults() {
	Default.URL = "https://default"
}

func setLogger() {}
