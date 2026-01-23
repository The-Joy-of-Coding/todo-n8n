package config

import (
	"strings"
	"testing"
)

func TestApiUrl(t *testing.T) {
	// This needs to test if there is a function called ApiUrl.
	// then check if it returns a string or not.
	// then check if the string is empty or not
	// then check if the string is in https format or not to pass this test.
	result := ApiUrl()

	if result == "" {
		t.Errorf("ApiUrl returned an empty string!")
	}

	if !strings.HasPrefix(result, "https://") {
		t.Errorf("ApiUrl() = %q; want it to start with 'https://'", result)
	}
}
