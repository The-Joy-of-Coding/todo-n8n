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
	t.Task = strings.Join(strings.Fields(t.Task), " ")
	if t.Task == "" {
		return
	}
	task := []rune(t.Task)
	task[0] = unicode.ToUpper(task[0])
	last := task[len(task)-1]
	if !unicode.IsPunct(last) {
		task = append(task, '.')
	}
	t.Task = string(task)
}
