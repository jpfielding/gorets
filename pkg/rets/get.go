package rets

import (
	"io/ioutil"
	"net/http"

	"context"
)

// GetRequest ...
type GetRequest struct {
	URL string
}

// Get gets an arbitrary file from the server or performs an arbitrary action, specified by URI
func Get(ctx context.Context, requester Requester, r GetRequest) error {
	req, err := http.NewRequest("GET", r.URL, nil)
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
