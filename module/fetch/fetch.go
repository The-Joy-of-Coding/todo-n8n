package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"todo-n8n/config"
	"todo-n8n/module/types"
)

type Transport types.Transport

var (
	client http.Client = http.Client{
		Timeout: config.GetTimeout(),
	}
	header string
	key    string
)

func init() {
	var err error
	header, err = config.GetEnv("AUTH_HEADER")
	if err != nil {
		slog.Error(err.Error())
		header = "HEADER"
	}
	key, err = config.GetEnv("AUTH_KEY")
	if err != nil {
		slog.Error(err.Error())
		key = "KEY"
	}
}

func (t *Transport) createRequest(url string, method string, body any) *Transport {
	if t.Error != nil {
		t = &Transport{}
		return t
	}
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			slog.Error(err.Error())
			t.Error = err
			return t
		}
	}
	req, err := http.NewRequest(
		method, config.GetURL(url),
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		slog.Error(err.Error())
		t.Error = err
		return t
	}
	t.Request = req
	return t
}

// Need to add a queue system and a temp cache to improve performance
func (t *Transport) fetch() *Transport {
	if t.Request == nil {
		slog.Info("The Request is empty!")
		return nil
	}
	t.Request.Header.Add("Content-Type", "application/json")
	t.Request.Header.Add(header, key)
	resp, err := client.Do(t.Request)
	if err != nil {
		slog.Error(err.Error())
		t.Error = err
		return t
	}
	t.Response = resp
	return t
}

func (t *Transport) validate(target any) error {
	if target == nil {
		t.Error = fmt.Errorf("Target is empty!")
	}
	if t.Response == nil || t.Response.Body == nil {
		t.Error = fmt.Errorf("no response body available to parse")
	}
	if t.Response != nil && t.Response.StatusCode >= 400 {
		t.Error = fmt.Errorf("server returned error: %d", t.Response.StatusCode)
	}
	return t.parseData(target)
}

func (t *Transport) parseData(target any) error {
	if t.Response != nil {
		defer t.Response.Body.Close()
	}
	if t.Error != nil {
		slog.Error(t.Error.Error())
		return t.Error
	}
	bodyBytes, err := io.ReadAll(t.Response.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	if string(bodyBytes) == "" {
		return fmt.Errorf("Server Failed to send a response!")
	}
	if err := json.Unmarshal(bodyBytes, target); err != nil {
		slog.Error("JSON unmarshal failed", "err", err, "body", string(bodyBytes))
		return err
	}
	return nil
}

func (t *Transport) processPendingReq() {}
