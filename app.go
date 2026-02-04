package main

import (
	"flag"
	"log/slog"

	"todo-n8n/module"
)

func main() {
	var isCommand bool
	getTasks := flag.Bool("get", false, "Returns all the pending todos")
	addTask := flag.String("add", "", "Give task details to add to todo")
	deleteTask := flag.Int("delete", 0, "Give task id to delete the todo")
	flag.Parse()
	if *addTask != "" {
		module.AddTask(*addTask)
		isCommand = true
	}
	if *deleteTask != 0 {
		module.DeleteTask(*deleteTask)
		isCommand = true
	}
	if *getTasks {
		module.GetTodos()
		isCommand = true
	}
	if isCommand {
		return
	}
	if err := module.Default(); err != nil {
		slog.Error(err.Error())
	}
}
