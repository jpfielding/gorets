/**
	extraction of the data pieces describing a RETS system
 */
package gorets


import (
	"fmt"
	"io/ioutil"
	"net/http"
)


type Metadata struct {

}

func (s *Session) GetMetadata(url, format, id, mtype string) (*Metadata, error) {
	qs := fmt.Sprintf("Format=%s",format)
	if id != "" {
		qs = qs +"&"+ fmt.Sprintf("ID=%s",id)
	}
	if mtype != "" {
		qs = qs +"&"+ fmt.Sprintf("Type=%s",mtype)
	}

	req, err := http.NewRequest(s.HttpMethod, url+"?"+qs, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	defer resp.Body.Close()

	switch mtype {
	case "METADATA-SYSTEM":
	case "METADATA-RESOURCE":

	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	metadata := &Metadata{}
	return metadata, nil
}
