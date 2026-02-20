package main

import (
	"flag"
	"log/slog"

	"todo-n8n/config"
	"todo-n8n/module"
)

func main() {
	config.SetLogger(false)
	isCommand := getArgs()
	if isCommand {
		return
	}
	if err := module.Default(); err != nil {
		slog.Error(err.Error())
	}
}

func getArgs() bool {
	var (
		err        error
		listTasks  = flag.Bool("list", false, "List all pending tasks")
		addTask    = flag.String("add", "", "Add a new task description")
		checkTask  = flag.Int("id", -1, "The ID of the task to act upon")
		updateTask = flag.String("update", "", "New text for the task")
		deleteTask = flag.Int("delete", -1, "Delete the task by ID")
	)
	flag.Parse()
	switch {
	case *addTask != "":
		err = module.AddTask(*addTask)
	case *deleteTask != -1:
		err = module.DeleteTask(*deleteTask)
	case *checkTask != -1:
		err = module.UpdateTask(*checkTask, *updateTask)
	case *listTasks:
		err = module.GetTodos()
	default:
		return false
	}
	if err != nil {
		slog.Error(err.Error())
	}
	return true
}
