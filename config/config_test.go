package config

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func TestPing(t *testing.T) {
	type Test struct {
		name     string
		inputUrl string
		result   bool
		matchLog string
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "error") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	tests := []Test{
		{"EMPTY URL", "", false, "URL is empty, please provide a proper url!"},
		{"INVALID URL", server.URL + "/error", false, "URL is invalid, please check the url provided."},
		{"NETWORK ERROR", "ftp://not-http", false, "Network error occured, please try again later or check again!"},
		{"VALID URL", server.URL, true, ""},
	}
	testFunc := func(test Test, t *testing.T) {
		var buf bytes.Buffer
		l := slog.New(slog.NewJSONHandler(&buf, nil))
		res, statusCode := ping(l, test.inputUrl)
		if !strings.Contains(buf.String(), test.matchLog) {
			t.Errorf("Missing log \"%s\" to indicate user!", test.matchLog)
		}
		if res != test.result {
			t.Errorf("Result mismatch: got %v, want %v please check your code!", res, test.result)
		}
		if statusCode >= 300 && res {
			t.Error("Test returned true for a non-ok status code")
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFunc(test, t)
		})
	}
}

func TestDefaults(t *testing.T) {
	if Default.URL == "" {
		t.Errorf("expected URL to be populated, got %q", Default.URL)
	}
	if Default.Timeout == 0 {
		t.Errorf("expected Timeout to be non-zero, got %v", Default.Timeout)
	}
}
