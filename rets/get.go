package rets

import (
	"io/ioutil"
	"net/http"

	"context"
)

// GetRequest ...
type GetRequest struct {
	URL, HTTPMethod string
}

// Get ...
func Get(requester Requester, ctx context.Context, r GetRequest) error {
	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return err
	}

	resp, err := requester(ctx, req)
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
