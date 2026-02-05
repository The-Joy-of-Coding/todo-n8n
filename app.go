package main

import (
	"flag"
	"log/slog"

	"todo-n8n/module"
)

func main() {
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
		listTasks  = flag.Bool("list", false, "List all pending tasks")
		addTask    = flag.String("add", "", "Add a new task description")
		checkTask  = flag.Int("id", -1, "The ID of the task to act upon")
		updateTask = flag.String("update", "", "New text for the task")
		deleteTask = flag.Int("delete", -1, "Delete the task by ID")
	)
	flag.Parse()
	switch {
	case *addTask != "":
		module.AddTask(*addTask)
	case *deleteTask != -1:
		module.DeleteTask(*deleteTask)
	case *checkTask != -1:
		module.UpdateTask(*checkTask, *updateTask)
	case *listTasks:
		module.GetTodos()
	default:
		return false
	}
	return true
}
