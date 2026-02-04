package module

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

//go:embed app.sh
var app_ui string

type Todos struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
}

type N8nRespnce struct {
	TodoList []Todos `json:"todo_list"`
}

func GetSafeTodos() {
	todo := Todos{}
	_, err := todo.get()
	if err != nil {
		slog.Error(err.Error())
	}
}

func AddTask(task string) {
	todo := Todos{Task: task}
	if err := todo.post(); err != nil {
		slog.Error(err.Error())
	}
}

func UpdateTask(id int, task string) {
	todo := Todos{Id: id}
	if task != "" {
		todo.Task = task
	}
	if err := todo.put(); err != nil {
		slog.Error(err.Error())
	}
}

func DeleteTask(id int) {
	todo := Todos{Id: id}
	if err := todo.delete(); err != nil {
		slog.Error(err.Error())
	}
}

func runner(command string) error {
	script := fmt.Sprintf("%s %s", app_ui, command)
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

func Default() error {
	return runner("todo_start")
}
