package main

import (
	"flag"
	"log/slog"

	"todo-n8n/module"
)

func main() {
	addData := flag.String("add", "*/*-", "Give task details to add to todo")
	deleteData := flag.Int("delete", 0, "Give task id to delete the todo")
	flag.Parse()
	if *addData != "*/*-" && *addData != "" {
		slog.Info("Task", "details", *addData)
		return
	}
	if *deleteData != 0 {
		slog.Info("Id", "details", *deleteData)
		return
	}
	if err := module.Default(); err != nil {
		slog.Error(err.Error())
	}
}
