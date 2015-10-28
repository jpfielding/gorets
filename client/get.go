package client

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
)

// GetRequest ...
type GetRequest struct {
	URL, HTTPMethod string
}

// Get ...
func (s *Session) Get(ctx context.Context, r GetRequest) error {
	method := s.HTTPMethodDefault
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return err
	}

	resp, err := s.Execute(ctx, req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
