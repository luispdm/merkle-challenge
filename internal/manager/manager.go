package manager

import (
	"fmt"
	"net/http"
)

const (
	Piece  = "piece"
	Hashes = "hashes"
)

var errFmt = "HTTP '%s' to '/%s' returned status code '%d'"

type IHTTP interface {
	Get(path string) (*http.Response, error)
}

type Manager struct {
	d decoder
	c IHTTP
}

type Option func(*Manager) error

type Request struct {
	path string
}

func NewRequest(path string) *Request {
	return &Request{
		path: path,
	}
}

func New(c IHTTP, opts ...Option) (*Manager, error) {
	m := &Manager{
		c: c,
		d: newJayson(),
	}
	for _, o := range opts {
		err := o(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func WithDecoder(d decoder) Option {
	return func(m *Manager) error {
		if d == nil {
			return fmt.Errorf("manager.WithDecoder: decoder can't be nil")
		}
		m.d = d
		return nil
	}
}

func (m *Manager) Get(r *Request, respObj interface{}) (int, error) {
	res, err := m.c.Get(r.path)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	err = m.d.decode(res, &respObj)
	if err != nil {
		return res.StatusCode, fmt.Errorf(errFmt, http.MethodGet, r.path, res.StatusCode)
	}
	return res.StatusCode, nil
}
