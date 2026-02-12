package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"todo-n8n/config"
)

type Transport struct {
	request  *http.Request
	response *http.Response
	error    error
}

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
		os.Exit(1)
	}
	key, err = config.GetEnv("AUTH_KEY")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (t *Transport) createRequest(method string, body any) *Transport {
	if t.error != nil {
		return t
	}
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			slog.Error(err.Error())
			t.error = err
			return t
		}
	}
	req, err := http.NewRequest(
		method, config.GetURL("API_URL"),
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		slog.Error(err.Error())
		t.error = err
		return t
	}
	t.request = req
	return t
}

// Need to add a queue system and a temp cache to improve performance
func (t *Transport) fetch() *Transport {
	t.request.Header.Add("Content-Type", "application/json")
	t.request.Header.Add(header, key)
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
	if target == nil {
		return fmt.Errorf("Target is empty!")
	}
	bodyBytes, err := io.ReadAll(t.response.Body)
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
