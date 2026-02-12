package template

import (
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	script, err := GetTemplate()
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
