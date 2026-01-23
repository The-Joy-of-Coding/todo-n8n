package config

import (
	"strings"
	"testing"
)

func TestApiUrl(t *testing.T) {
	result := ApiUrl()

	if result == "" {
		t.Errorf("ApiUrl returned an empty string!")
	}

	if !strings.HasPrefix(result, "https://") {
		t.Errorf("ApiUrl() = %q; want it to start with 'https://'", result)
	}
}
