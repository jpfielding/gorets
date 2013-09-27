/**

 */
package gorets


import (
	"net/http"
)


type Metadata struct {

}


func (s *Session) GetMetadata(url string) (*Metadata, error) {
	// TODO setup request

	req, err := http.NewRequest(s.HttpMethod, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	metadata := &Metadata{}
	return metadata, nil
}
