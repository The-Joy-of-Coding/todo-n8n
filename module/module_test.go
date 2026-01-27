package module

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todo-n8n/config"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()
	buff := config.SetLogger(true)
	req, _ := http.NewRequest("GET", server.URL, nil)
	trans := Transport{request: req}
	trans.fetch()
	if trans.responce == nil {
		t.Fatalf("Fetch failed: response is nil. Logs: %s", buff.Buff.String())
	}
	if trans.responce.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", trans.responce.StatusCode)
	}
}

func TestParseData(t *testing.T) {
	type MockResponse struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}
	tests := []struct {
		name    string
		json    string
		nilRes  bool
		wantErr bool
	}{
		{"Valid JSON", `{"status": "success", "id": 123}`, false, false},
		{"Invalid JSON", `{"status": "success", "id": "not-an-int"}`, false, true},
		{"Nil Response", "", true, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var trans Transport
			if !tc.nilRes {
				res := &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(tc.json)),
				}
				trans.responce = res
			}
			var result MockResponse
			err := trans.ParseData(&result)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

