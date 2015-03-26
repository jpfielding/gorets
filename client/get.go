/**
Maybe the dumbest thing in all of RETS
*/
package client

import (
	"io/ioutil"
	"net/http"
)

type GetRequest struct {
	URL, HTTPMethod string
}

/**
TODO - this needs to somehow send the results back to the caller
*/
func (s *Session) Get(r GetRequest) error {
	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return err
	}

	resp, err := s.Client.Do(req)
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
