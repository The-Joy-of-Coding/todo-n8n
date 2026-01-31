package config

import (
	"bufio"
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
	SetLogger(false)
	if err := json.Unmarshal(file, &_default); err != nil {
		slog.Error("failed to unmarshal defaults", "error", err)
	}
	_ = LoadEnv(".env") || LoadEnv("../.env")
}

func LoadEnv(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
	return true
}

func GetEnv(env string) (string, error) {
	res := os.Getenv(env)
	if res == "" {
		return res, fmt.Errorf("The result is empty!")
	}
	return res, nil
}

func GetURL(key string) string {
	val, err := GetEnv(key)
	if err != nil {
		slog.Error("No valid url found!", "key", key)
		val, err = GetEnv("API_URL")
		if err != nil {
			return "http://localhost:5678"
		}
	}
	if val == "" {
		return "http://localhost:5678"
	}
	if !strings.HasPrefix(val, "https://") {
		slog.Warn("Using insecure api endpoint.")
	}
	return val
}

func ping(targetURL string) (bool, int) {
	if targetURL == "" {
		slog.Error("URL is empty, please provide a proper url!")
		return false, 0
	}
	client := http.Client{Timeout: GetTimeout()}
	resp, err := client.Get(targetURL)
	if err != nil {
		slog.Error("Network error occured, please try again later or check again!")
		return false, 0
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		slog.Error("URL is invalid, please check the url provided.", "status", resp.StatusCode)
		return false, resp.StatusCode
	}
	return true, resp.StatusCode
}

func GetTimeout() time.Duration {
	if _default.Timeout > 0 {
		return time.Duration(_default.Timeout) * time.Second
	}
	return 10 * time.Second
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
