package module

type Todos struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
}

var todos []Todos

func init() {
	todo := Todos{}
	todo.get()
}

func AddTask(task string) []Todos {
	todo := Todos{Task: task}
	todo.post()
	return todos
}

func CheckOrDeleteTask(id int, isCheck bool) []Todos {
	todo := Todos{Id: id}
	if isCheck {
		todo.put()
	} else {
		todo.delete()
	}
	return todos
}

func Default() {}
