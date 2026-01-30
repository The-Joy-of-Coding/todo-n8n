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

func TestTodoLifecycle(t *testing.T) {
	type Test struct {
		name     string
		taskName string
		isCheck  bool
		wantLog  string
	}
	var testTaskID int
	testFunc := func(tc Test, t *testing.T) {
		_ = config.SetLogger(true)
		switch tc.name {
		case "POST TASK":
			res := AddTask(tc.taskName)
			found := false
			for _, item := range res {
				if item.Task == tc.taskName {
					testTaskID = item.Id
					found = true
				}
			}
			if !found || testTaskID == 0 {
				t.Errorf("Failed to create or find task: %s", tc.taskName)
			}
		case "PUT TASK":
			if testTaskID == 0 {
				t.Skip("No ID available from POST")
			}
			res := CheckOrDeleteTask(testTaskID, true)
			found := false
			for _, item := range res {
				if item.Id == testTaskID {
					found = true
				}
			}
			if !found {
				t.Errorf("Task ID %d disappeared after PUT", testTaskID)
			}
		case "GET TASKS":
			todo := Todos{}
			todos, err := todo.get()
			if err != nil {
				t.Errorf("GET failed: %v", err)
			}
			if len(todos.TodoList) == 0 {
				t.Error("Todos slice is empty after GET")
			}
		case "DELETE TASK":
			if testTaskID == 0 {
				t.Skip("No ID available from POST")
			}
			res := CheckOrDeleteTask(testTaskID, false)
			for _, item := range res {
				if item.Id == testTaskID {
					t.Errorf("Task ID %d still exists after DELETE", testTaskID)
				}
			}
		}
	}
	tests := []Test{
		{name: "POST TASK", taskName: "N8N Integration Test"},
		{name: "PUT TASK", isCheck: true},
		{name: "GET TASKS"},
		{name: "DELETE TASK", isCheck: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testFunc(tc, t)
		})
	}
}
