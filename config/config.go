package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type configDefault struct {
	URLENV  string `json:"urlEnv"`
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
}

type Test struct {
	Buff *bytes.Buffer
}

var (
	//go:embed defaults.json
	file     []byte
	_default configDefault
)

func init() {
	if err := json.Unmarshal(file, &_default); err != nil {
		slog.Error(err.Error())
	}
}

func GetEnv(env string) (string, error) {
	res := os.Getenv(env)
	if res == "" {
		return res, fmt.Errorf("The result is empty!")
	}
	return res, nil
}

func ApiUrl(url_key string) string {
	url := os.Getenv(url_key)
	if url == "" {
		slog.Error("No valid url found!")
		return _default.Url()
	}
	if !strings.HasPrefix(url, "https://") {
		slog.Warn("Using insecure api endpoint.")
	}
	if ok, _ := ping(url); !ok {
		slog.Error("Url is invalid, please change the url to run the application.")
		return ApiUrl(_default.URLENV)
	}
	return url
}

func ping(url string) (bool, int) {
	if url == "" {
		slog.Error("URL is empty, please provide a proper url!")
		return false, 0
	}
	client := http.Client{
		Timeout: TimeOut(),
	}
	resp, err := client.Head(url)
	if err != nil {
		slog.Error("Network error occured, please try again later or check again!")
		return false, 0
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		slog.Error("URL is invalid, please check the url provided.")
		return false, http.StatusNotFound
	}
	return true, resp.StatusCode
}

func TimeOut() time.Duration {
	timeOut := 5 * time.Second
	if _default.Timeout > 0 {
		timeOut = time.Duration(_default.Timeout) * time.Second
	}
	return timeOut
}

func (c *configDefault) Url() string {
	if c.URL == "" {
		return "http://localhost:5678"
	}
	if ok, _ := ping(_default.URL); !ok {
		return "http://localhost:5678"
	}
	return c.URL
}

func SetLogger(test bool) *Test {
	var buff bytes.Buffer
	var handler slog.Handler
	if test {
		handler = slog.NewJSONHandler(&buff, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	}
	slog.SetDefault(
		slog.New(handler),
	)
	if test {
		return &Test{
			Buff: &buff,
		}
	}
	return nil
}
