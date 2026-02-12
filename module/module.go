package module

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"time"

	"todo-n8n/module/fetch"
	"todo-n8n/module/template"
)

func GetTodos() error {
	todo := fetch.Todos{}
	data, err := todo.Get()
	for _, v := range data.TodoList {
		fmt.Printf("Id: %v - Task: %s\n", v.Id, v.Task)
	}
	return err
}

func AddTask(task string) error {
	todo := fetch.Todos{Task: task}
	err := todo.Post()
	return err
}

func UpdateTask(id int, task string) error {
	todo := fetch.Todos{Id: id}
	if task != "" {
		todo.Task = task
	}
	err := todo.Put()
	return err
}

func DeleteTask(id int) error {
	todo := fetch.Todos{Id: id}
	err := todo.Delete()
	return err
}

func Default() error {
	app_ui, err := template.GetTemplate()
	if err != nil {
		return err
	}
	script := fmt.Sprintf("%s todo_start", app_ui)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Minute*5,
	)
	defer cancel()
	cmd := exec.CommandContext(
		ctx, "bash", "-c", script,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
