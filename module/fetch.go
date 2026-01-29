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
	response *http.Response
	error    error
}

func (t *Todos) get() {
	req, err := http.NewRequest("GET", config.ApiUrl(""), nil)
	if err != nil {
		slog.Error(err.Error())
	}
	tr := Transport{request: req}
	tr.fetch().ParseData(&todos)
}

func (t *Todos) post() {
	jsonBody, err := json.Marshal(t)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"POST", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.get()
}

func (t *Todos) put() {
	jsonBody, err := json.Marshal(t)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"PUT", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.get()
}

func (t *Todos) delete() {
	jsonBody, err := json.Marshal(t)
	if err != nil {
		slog.Error(err.Error())
	}
	req, err := http.NewRequest(
		"DELETE", config.ApiUrl(""),
		bytes.NewReader(jsonBody),
	)
	tr := Transport{request: req}
	tr.fetch().ParseData(nil)
	t.get()
}

func (t *Transport) fetch() *Transport {
	client := http.Client{
		Timeout: config.TimeOut(),
	}
	if t.request.Body != nil {
		t.request.Header.Add("Content-Type", "application/json")
	}
	// t.request.Header.Add("", "")
	resp, err := client.Do(t.request)
	if err != nil {
		slog.Error(err.Error())
		t.error = err
		return t
	}
	t.response = resp
	return t
}

func (t *Transport) ParseData(target any) error {
	if t.error != nil {
		return t.error
	}
	if t.response == nil || t.response.Body == nil {
		err := "no response body available to parse"
		slog.Error(err)
		return fmt.Errorf("%s", err)
	}
	defer t.response.Body.Close()
	if t.response.StatusCode >= 400 {
		err := fmt.Errorf("server returned error: %d", t.response.StatusCode)
		slog.Error(err.Error())
		return err
	}

	return json.NewDecoder(t.response.Body).Decode(target)
}
