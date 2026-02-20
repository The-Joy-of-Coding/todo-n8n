package fetch

import (
	"fmt"

	"todo-n8n/module/types"
)

type Todos types.Todos

func (t *Todos) Get() (types.N8nResponse, error) {
	var n8n types.N8nResponse
	tr := Transport{}
	return n8n, tr.
		createRequest("", "GET", nil).
		fetch().
		validate(&n8n)
}

func (t *Todos) Post() error {
	tr := Transport{}
	var res any
	err := tr.
		createRequest("", "POST", t).
		fetch().
		validate(&res)
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
		createRequest("", "PUT", t).
		fetch().
		validate(&res)
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
		createRequest("", "DELETE", t).
		fetch().
		validate(&res)
	if err != nil {
		return err
	}
	if res.Task == "" {
		return fmt.Errorf("Something went wrong! please check the result: %v", res)
	}
	return nil
}
