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
	if trans.response == nil {
		t.Fatalf("Fetch failed: response is nil. Logs: %s", buff.Buff.String())
	}
	if trans.response.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", trans.response.StatusCode)
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
				trans.response = res
			}
			var result MockResponse
			err := trans.ParseData(&result)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	testCase := Todos{Task: "This is a test request from go."}
	t.Run("POST_RESQUEST", func(t *testing.T) {
		todo := Todos{
			Task: testCase.Task,
		}
		err := todo.post()
		if err != nil {
			t.Errorf("POST request test failed due to: %s", err.Error())
		}
	})
	t.Run("GET_REQUEST", func(t *testing.T) {
		todo := Todos{}
		res, err := todo.get()
		if err != nil {
			t.Errorf("GET requst test failed due to: %s", err.Error())
		}
		for _, v := range res.TodoList {
			if v.Task == testCase.Task {
				testCase.Id = v.Id
			}
		}
		if testCase.Id == 0 {
			t.Errorf("GET request test failed due to: unable to find the test-case task")
		}
	})
	t.Run("PUT_REQUEST", func(t *testing.T) {
		todo := Todos{Id: testCase.Id}
		err := todo.put()
		if err != nil {
			t.Errorf("PUT request test failed due to: %s", err.Error())
		}
	})
	t.Run("DELETE_REQUEST", func(t *testing.T) {
		todo := Todos{Id: testCase.Id}
		err := todo.delete()
		if err != nil {
			t.Errorf("DELETE request test failed due to: %s", err.Error())
		}
	})
}

func TestGetTemplate(t *testing.T) {
	script, err := getTemplate()
	if err != nil {
		t.Fatalf("getTemplate failed: %v", err)
	}
	if len(script) == 0 {
		t.Error("expected non-empty script, got empty string")
	}
	expectedFunctions := []string{
		"todo_header()",
		"todo_menu()",
		"todo_start()",
	}
	for _, fn := range expectedFunctions {
		if !strings.Contains(script, fn) {
			t.Errorf("script missing expected function: %s", fn)
		}
	}
	if strings.Contains(script, "random_text_that_shouldnt_be_there") {
		t.Error("script contains unexpected content")
	}
}
