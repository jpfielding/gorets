/**
	Maybe the dumbest thing in all of RETS
 */
package gorets

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/**
	no, really, thats all it does, i probably need to parse some pointless response
 */
func (s *Session) Get(url string) (error) {
	req, err := http.NewRequest(s.HttpMethod, url, nil)
	if err != nil {
		return err
	}

	resp, err := s.Client.Do(req)
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}