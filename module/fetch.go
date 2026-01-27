package module

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"todo-n8n/config"
)

type Transport struct {
	request  *http.Request
	responce *http.Response
}

func (t *todos) Get() {
	req, err := http.NewRequest("GET", config.ApiUrl(""), nil)
	if err != nil {
		slog.Error(err.Error())
	}
	tr := Transport{request: req}
	tr.fetch().ParseData(t)
}

func (t *todos) Post() {
	t.Get()
} // Post will add a new todo to the todos via api.

func (t *todos) Put() {
	t.Get()
} // Put will mark the required todo as checked via api.

func (t *todos) Delete() {
	t.Get()
} // Delete will delete the todo from todos via api.

func (t *Transport) fetch() *Transport {
	client := http.Client{
		Timeout: config.TimeOut(),
	}
	resp, err := client.Do(t.request)
	if err != nil {
		slog.Error(err.Error())
		return t
	}
	t.responce = resp
	return t
}

func (t *Transport) ParseData(target interface{}) error {
	if t.responce == nil || t.responce.Body == nil {
		err := "no response body available to parse"
		slog.Error(err)
		return fmt.Errorf("%s", err)
	}
	defer t.responce.Body.Close()
	if err := json.NewDecoder(t.responce.Body).Decode(target); err != nil {
		slog.Error("json decode failed", "error", err)
		return err
	}
	return nil
}
