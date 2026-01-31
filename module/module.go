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

func CheckTask(id int) {
	todo := Todos{Id: id}
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

func runner(script string) error {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Minute*5,
	)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Default() error {
	script := fmt.Sprintf("%s todo_start", app_ui)
	slog.Info(script)
	return runner(script)
}
