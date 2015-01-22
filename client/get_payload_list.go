package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type PayloadList struct {
	Rets     RetsResponse
	Error    error
	Payloads <-chan CompactData
}

type PayloadListRequest struct {
	Url, Id string
}

/**
 */
func (s *Session) GetPayloadList(p PayloadListRequest) (*PayloadList, error) {
	// required
	values := url.Values{}
	if p.Id == "" {
		values.Add("ID", p.Id)
	}

	req, err := http.NewRequest(s.HttpMethod, p.Url+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseGetPayloadList(resp.Body)
}

func parseGetPayloadList(body io.ReadCloser) (*PayloadList, error) {
	payloads := make(chan CompactData)
	parser := xml.NewDecoder(body)

	list := PayloadList{
		Payloads: payloads,
	}
	delim := "	"
	// backgroundable processing of the data into our buffer
	dataProcessing := func() {
		// this channel needs to be closed or the caller can infinite loop
		defer close(payloads)
		// this is the web socket that needs addressed
		defer body.Close()
		// extract the data
		for {
			token, err := parser.Token()
			if err != nil {
				list.Error = err
				break
			}
			switch t := token.(type) {
			case xml.StartElement:
				switch t.Name.Local {
				case "RETSPayloadList":
					mcd, err := ParseMetadataCompactDecoded(t, parser, delim)
					if err != nil {
						fmt.Println("failed to decode: ", err)
						list.Error = err
						return
					}
					payloads <- *mcd
				}
			}
		}
	}
	for {
		token, err := parser.Token()
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			switch elmt.Name.Local {
			case "RETS", "RETS-STATUS":
				rets, err := ParseRetsResponseTag(elmt)
				if err != nil {
					return nil, err
				}
				list.Rets = *rets
				go dataProcessing()
				return &list, nil
			case "DELIMITER":
				decoded, err := ParseDelimiterTag(elmt)
				if err != nil {
					return nil, err
				}
				delim = decoded
			}
		}
	}
}
