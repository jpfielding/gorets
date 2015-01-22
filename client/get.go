/**
Maybe the dumbest thing in all of RETS
*/
package client

import (
	"io/ioutil"
	"net/http"
)

/**
TODO - this needs to somehow send the results back to the caller
*/
func (s *Session) Get(url string) error {
	req, err := http.NewRequest(s.HttpMethod, url, nil)
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
