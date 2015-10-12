/**
Maybe the dumbest thing in all of RETS
*/
package client

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type GetRequest struct {
	URL, HTTPMethod string
}

/**
TODO - this needs to somehow send the results back to the caller
*/
func (s *Session) Get(ctx context.Context, r GetRequest) error {
	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return err
	}

	resp, err := ctxhttp.Do(ctx, &s.Client, req)
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
