package fetch

import (
	"fmt"
	"strings"
	"unicode"
)

func (t *Todos) Validate() error {
	words := strings.Fields(t.Task)
	count := len(words)
	minCount, maxCount := 5, 50
	switch {
	case count < minCount:
		return fmt.Errorf("The task is shorter then 5 words")
	case count > maxCount:
		return fmt.Errorf("The task is longer the 50 words")
	default:
		return nil
	}
}

func (t *Todos) Formate() {
	t.Task = strings.TrimSpace(t.Task)
	if len(t.Task) == 0 {
		return
	}
	task := []rune(t.Task)
	task[0] = unicode.ToUpper(task[0])
	t.Task = string(task)
}
