package module

import (
	"log/slog"
)

type Todos struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
}

type N8nRespnce struct {
	TodoList []Todos `json:"todo_list"`
}

func getSafeTodos() []Todos {
	fail := make(chan N8nRespnce)
	go func(fail chan N8nRespnce) {
		todo := Todos{}
		res, err := todo.get()
		if err != nil {
			slog.Error(err.Error())
			fail <- res
			return
		}
		fail <- res
	}(fail)
	res := <-fail
	if len(res.TodoList) == 0 {
		return nil
	}
	slog.Info("Todos", "List", res.TodoList)
	return res.TodoList
}

func AddTask(task string) []Todos {
	todo := Todos{Task: task}
	if err := todo.post(); err != nil {
		slog.Error(err.Error())
		return getSafeTodos()
	}
	return getSafeTodos()
}

func CheckOrDeleteTask(id int, isCheck bool) []Todos {
	todo := Todos{Id: id}
	if isCheck {
		todo.put()
	} else {
		todo.delete()
	}
	return getSafeTodos()
}

func Default() {}
