package module

import (
	"bytes"
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

func (t *todos) Post(task string) {
	data := map[string]string{"data": task}
	jsonBody, err := json.Marshal(data)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"POST", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.Get()
}

func (t *todos) Put(id int) {
	data := map[string]int{"id": id}
	jsonBody, err := json.Marshal(data)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"PUT", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.Get()
}

func (t *todos) Delete(id int) {
	data := map[string]int{"id": id}
	jsonBody, err := json.Marshal(data)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"DELETE", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.Get()
}

func (t *Transport) fetch() *Transport {
	client := http.Client{
		Timeout: config.TimeOut(),
	}
	t.request.Header.Add("", "")
	resp, err := client.Do(t.request)
	if err != nil {
		slog.Error(err.Error())
		return t
	}
	t.responce = resp
	return t
}

func (t *Transport) ParseData(target any) error {
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
