package types

import "net/http"

type Transport struct {
	Request  *http.Request
	Response *http.Response
	Error    error
}
