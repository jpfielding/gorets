/**
     notify the server we're done and kill our session
*/
package gorets

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

func (s *Session) Logout(logoutUrl string) error {
	req, err := http.NewRequest(s.HttpMethod, logoutUrl, nil)
	if err != nil {
		return err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// wipe the cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.Client.Jar = jar
	return nil
}
