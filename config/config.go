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
	file     []byte
	_default ConfigDefault
)

func init() {
	if err := json.Unmarshal(file, &_default); err != nil {
		slog.Error(err.Error())
	}
}

func ApiUrl(logger *slog.Logger, url_key string) string {
	url := os.Getenv(url_key)
	if url == "" {
		logger.Error("No valid url found!")
		return _default.Url()
	}
	if !strings.HasPrefix(url, "https://") {
		logger.Warn("Using insecure api endpoint.")
	}
	if ok, _ := ping(logger, url); !ok {
		logger.Error("Url is invalid, please change the url to run the application.")
		return _default.Url()
	}
	return url
}

func ping(logger *slog.Logger, url string) (bool, int) {
	if url == "" {
		logger.Error("URL is empty, please provide a proper url!")
		return false, 0
	}
	client := http.Client{
		Timeout: _default.timeOut(),
	}
	resp, err := client.Head(url)
	if err != nil {
		logger.Error("Network error occured, please try again later or check again!")
		return false, resp.StatusCode
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		logger.Error("URL is invalid, please check the url provided.")
		return false, http.StatusNotFound
	}
	return true, resp.StatusCode
}

func (c *ConfigDefault) timeOut() time.Duration {
	timeOut := 5 * time.Second
	if _default.Timeout > 0 {
		timeOut = time.Duration(c.Timeout) * time.Second
	}
	return timeOut
}

func (c *ConfigDefault) Url() string {
	if c.URL == "" {
		return "http://localhost:5678"
	}
	if ok, _ := ping(slog.Default(), c.URL); !ok {
		return "http://localhost:5678"
	}
	return c.URL
}
