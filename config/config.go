package config

import "os"

func ApiUrl() string {
	return os.Getenv("API_URL")
}
