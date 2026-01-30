package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
		buff := SetLogger(true)
		res := GetURL(test.inputUrl)
		if !strings.Contains(buff.Buff.String(), test.matchLog) {
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
		buff := SetLogger(true)
		res, statusCode := ping(test.inputUrl)
		if !strings.Contains(buff.Buff.String(), test.matchLog) {
			t.Errorf("Missing log \"%s\" to indicate user!", test.matchLog)
		}
		if res != test.result {
			t.Errorf("Result mismatch: got %v, want %v please check your code!", res, test.result)
		}
		if (statusCode >= 300 || statusCode < 200) && res {
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
	var url string = GetURL("")
	if url == "" {
		t.Errorf("expected URL to be populated, got %q", url)
	}
	var timeout int = int(GetTimeout())
	if timeout <= 0 {
		t.Errorf("expected Timeout to be non-zero, got %v", timeout)
	}
}

func TestEnv(t *testing.T) {
	type Test struct {
		name     string
		inputEnv string
		outPut   bool
	}
	tests := []Test{
		{"ENV EMPTY", "test", false},
		{"ENV NOT EMPTY", "INSECURE_API", true},
	}
	testFunc := func(test Test, t *testing.T) {
		res, err := GetEnv(test.inputEnv)
		if (err != nil) && test.outPut {
			t.Errorf("Somthing went wrong in your code, please check the GetEnv function, err: %s", err.Error())
		}
		if (res != "") != test.outPut {
			t.Errorf("Test out put is not same, please check your logic!")
		}
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testFunc(tc, t)
		})
	}
}

func TestLoadEnv(t *testing.T) {
	tmpDir := t.TempDir()
	type Test struct {
		name     string
		filename string
		content  string
		wantErr  bool
		checkKey string
		wantVal  string
	}
	tests := []Test{
		{
			name:     "VALID FILE",
			filename: "valid.env",
			content:  "TEST_KEY=hello_world\nANOTHER_VAR=123",
			wantErr:  false,
			checkKey: "TEST_KEY",
			wantVal:  "hello_world",
		},
		{
			name:     "NOISY FILE",
			filename: "noisy.env",
			content:  "\n# comment\n  SPACED_KEY = value  \n\n",
			wantErr:  false,
			checkKey: "SPACED_KEY",
			wantVal:  "value",
		},
		{
			name:     "MISSING FILE",
			filename: "missing.env",
			content:  "",
			wantErr:  true,
		},
	}
	testFunc := func(tc Test, t *testing.T) {
		path := filepath.Join(tmpDir, tc.filename)
		if tc.name != "MISSING FILE" {
			err := os.WriteFile(path, []byte(tc.content), 0o644)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
		}
		ok := LoadEnv(path)
		if !ok != tc.wantErr {
			t.Errorf("LoadEnv() error = %v, wantErr %v", ok, tc.wantErr)
			return
		}
		if tc.checkKey != "" {
			gotVal := os.Getenv(tc.checkKey)
			if gotVal != tc.wantVal {
				t.Errorf("Env mismatch: got %q, want %q", gotVal, tc.wantVal)
			}
			os.Unsetenv(tc.checkKey)
		}
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testFunc(tc, t)
		})
	}
}
