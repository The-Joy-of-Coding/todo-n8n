package config

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestApiUrl(t *testing.T) {
	type Test struct {
		name     string
		inputUrl string
		matchLog string
	}
	tests := []Test{
		{"EMPTY URL", "", "No valid url found!"},
		{"INSECURE URL", "INSECURE_API", "Using insecure api endpoint."},
		{"SECURE URL", "SECURE_API", ""},
	}
	testFunc := func(test Test, t *testing.T) {
		var buf bytes.Buffer
		l := slog.New(slog.NewJSONHandler(&buf, nil))
		res := ApiUrl(l, test.inputUrl)
		if !strings.Contains(buf.String(), test.matchLog) {
			t.Errorf("Missing log \"%s\" to indicate user!", test.matchLog)
		}
		if res == "" {
			t.Errorf("Returning empty result!")
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFunc(test, t)
		})
	}
}
