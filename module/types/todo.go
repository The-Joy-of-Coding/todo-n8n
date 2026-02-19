package types

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
