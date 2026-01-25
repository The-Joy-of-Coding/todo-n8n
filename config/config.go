package config

import (
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type ConfigDefault struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
}

var (
	//go:embed defaults.json
	file    []byte
	Default ConfigDefault
)

func init() {
	setDefaults()
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
	if ok, _ := ping(logger, url); !ok {
		logger.Error("Url is invalid, please change the url to run the application.")
		return Default.URL
	}
	Default.URL = url
	return url
}

func ping(logger *slog.Logger, url string) (bool, int) {
	if url == "" {
		logger.Error("URL is empty, please provide a proper url!")
		return false, 0
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		logger.Error("Network error occured, please try again later or check again!")
		return false, 0
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		logger.Error("URL is invalid, please check the url provided.")
		return false, 404
	}
	return true, resp.StatusCode
}

func setDefaults() {
	if err := json.Unmarshal(file, &Default); err != nil {
		slog.Error(err.Error())
	}
}
