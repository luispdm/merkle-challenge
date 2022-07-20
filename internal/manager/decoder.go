package manager

import (
	"encoding/json"
	"net/http"
)

type decoder interface {
	decode(*http.Response, interface{}) error
}

type jayson struct{}

func newJayson() *jayson {
	return &jayson{}
}

func (j *jayson) decode(res *http.Response, toFill interface{}) error {
	err := json.NewDecoder(res.Body).Decode(&toFill)
	if err != nil {
		return err
	}
	return nil
}
