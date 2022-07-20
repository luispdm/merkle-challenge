package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type HTTP struct {
	serverURL string
	client    *http.Client
}

func NewHTTP(sURL string, timeout time.Duration) *HTTP {
	return &HTTP{
		client:    &http.Client{Timeout: timeout},
		serverURL: sURL,
	}
}

func (c *HTTP) Get(path string) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, path)
}

func (c *HTTP) makeRequest(method, path string, body ...[]byte) (*http.Response, error) {
	b := new(bytes.Buffer)
	if len(body) > 0 {
		b = bytes.NewBuffer(body[0])
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.serverURL, path), b)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
