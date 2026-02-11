package module

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Todos struct {
	Id          int    `json:"id"`
	Task        string `json:"task"`
	Priority    int    `json:"priority"`
	ListOfShame bool   `json:"list_of_shame"`
	Deadline    string `json:"deadline"`
	HallOfFame  bool   `json:"hall_of_fame"`
}

type N8nResponse struct {
	TodoList []Todos `json:"todo_list"`
}

func GetTodos() error {
	todo := Todos{}
	data, err := todo.get()
	for _, v := range data.TodoList {
		fmt.Printf("Id: %v - Task: %s\n", v.Id, v.Task)
	}
	return err
}

func AddTask(task string) error {
	todo := Todos{Task: task}
	err := todo.post()
	return err
}

func UpdateTask(id int, task string) error {
	todo := Todos{Id: id}
	if task != "" {
		todo.Task = task
	}
	err := todo.put()
	return err
}

func DeleteTask(id int) error {
	todo := Todos{Id: id}
	err := todo.delete()
	return err
}

func Default() error {
	app_ui, err := getTemplate()
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
