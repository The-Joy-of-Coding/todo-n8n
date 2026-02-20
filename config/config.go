package config

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type configDefault struct {
	URL_ENV     string
	DEFAULT_URL string
	Timeout     int
}

type Test struct {
	Buff *bytes.Buffer
}

var _default = configDefault{
	DEFAULT_URL: "http://localhost:5678",
	URL_ENV:     "API_URL",
	Timeout:     10,
}

func GetEnv(env string) (string, error) {
	res := os.Getenv(env)
	if res == "" {
		return res, fmt.Errorf("The result is empty!")
	}
	return res, nil
}

func (c *configDefault) logError(log string, args ...any) string {
	slog.Error(log, args...)
	return c.DEFAULT_URL
}

func GetURL(key string) string {
	val, err := GetEnv(key)
	if err != nil || val == "" {
		slog.Warn("Given Api is empty!")
		return _default.logError("API url is empty, using the default url!")
	}
	if !strings.HasPrefix(val, "https://") {
		slog.Warn("Using insecure api endpoint.")
	}
	return val
}

func GetTimeout() time.Duration {
	return time.Duration(_default.Timeout) * time.Second
}

func SetLogger(test bool) *Test {
	var buff bytes.Buffer
	var writer io.Writer
	if test {
		writer = &buff
	} else {
		writer = os.Stderr
	}
	handler := slog.NewJSONHandler(
		writer, nil,
	)
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
