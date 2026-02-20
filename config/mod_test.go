package config

import (
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
		{"EMPTY URL", "", "Given Api is empty!"},
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
