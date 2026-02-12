package fetch

import (
	"fmt"
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

func (t *Todos) Get() (N8nResponse, error) {
	var n8n N8nResponse
	tr := Transport{}
	return n8n, tr.
		createRequest("GET", nil).
		fetch().
		ParseData(&n8n)
}

func (t *Todos) Post() error {
	tr := Transport{}
	var res any
	err := tr.
		createRequest("POST", t).
		fetch().
		ParseData(&res)
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("Something went wrong! please check the result: %v", res)
	}
	return nil
}

func (t *Todos) Put() error {
	tr := Transport{}
	var res Todos
	err := tr.
		createRequest("PUT", t).
		fetch().
		ParseData(&res)
	if err != nil {
		return err
	}
	if res.Task == "" {
		return fmt.Errorf("Something went wrong! please check the result: %v", res)
	}
	return nil
}

func (t *Todos) Delete() error {
	tr := Transport{}
	var res Todos
	err := tr.
		createRequest("DELETE", t).
		fetch().
		ParseData(&res)
	if err != nil {
		return err
	}
	if res.Task == "" {
		return fmt.Errorf("Something went wrong! please check the result: %v", res)
	}
	return nil
}
