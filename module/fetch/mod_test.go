package fetch

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
	trans := Transport{Request: req}
	trans.fetch()
	if trans.Response == nil {
		t.Fatalf("Fetch failed: response is nil. Logs: %s", buff.Buff.String())
	}
	if trans.Response.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", trans.Response.StatusCode)
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
				trans.Response = res
			}
			var result MockResponse
			err := trans.validate(&result)
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
		err := todo.Post()
		if err != nil {
			t.Errorf("POST request test failed due to: %s", err.Error())
		}
	})
	t.Run("GET_REQUEST", func(t *testing.T) {
		todo := Todos{}
		res, err := todo.Get()
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
		err := todo.Put()
		if err != nil {
			t.Errorf("PUT request test failed due to: %s", err.Error())
		}
	})
	t.Run("DELETE_REQUEST", func(t *testing.T) {
		todo := Todos{Id: testCase.Id}
		err := todo.Delete()
		if err != nil {
			t.Errorf("DELETE request test failed due to: %s", err.Error())
		}
	})
}

func TestValidate(t *testing.T) {
	type Tests struct {
		Name      string
		Input     Todos
		WantError bool
	}
	tests := []Tests{
		{
			Name: "Too short (4 words)",
			Input: Todos{
				Task: "This is too short.",
			},
			WantError: true,
		},
		{
			Name: "Minimum words (5 words)",
			Input: Todos{
				Task: "This task has five words.",
			},
			WantError: false,
		},
		{
			Name: "Valid length",
			Input: Todos{
				Task: "This is a task with more than five words but less than fifty.",
			},
			WantError: false,
		},
		{
			Name: "Too long (51 words)",
			Input: Todos{
				Task: "lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore magna aliqua ut enim ad minim veniam quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur",
			},
			WantError: true,
		},
	}
	testFunc := func(test Tests, t *testing.T) {
		err := test.Input.Validate()
		if (err != nil) != test.WantError {
			t.Errorf("Validate() error = %v, wantErr %v", err, test.WantError)
		}
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testFunc(test, t)
		})
	}
}

func TestFormate(t *testing.T) {
	type testStruct struct {
		Name     string
		Input    string
		Expected string
	}
	tests := []testStruct{
		{"Empty string", "", ""},
		{"Whitespace only", "   ", ""},
		{"Already formatted", "Task task", "Task task"},
		{"Lowercase start", "task task", "Task task"},
		{"With whitespace", "  task task  ", "Task task"},
	}
	testFunc := func(test testStruct, t *testing.T) {
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testFunc(test, t)
		})
	}
}
