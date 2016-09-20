package rets

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"context"
)

// TODO redo this without channels

// PayloadList ...
type PayloadList struct {
	Rets     RetsResponse
	Error    error
	Payloads <-chan CompactData
}

// PayloadListRequest ...
type PayloadListRequest struct {
	URL, HTTPMethod, ID string
}

// GetPayloadList ...
func GetPayloadList(requester Requester, ctx context.Context, p PayloadListRequest) (*PayloadList, error) {
	// required
	values := url.Values{}
	if p.ID == "" {
		values.Add("ID", p.ID)
	}

	method := "GET"
	if p.HTTPMethod != "" {
		method = p.HTTPMethod
	}
	// TODO use a URL object then properly append to it
	req, err := http.NewRequest(method, p.URL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
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
					cd, err := NewCompactData(t, parser, delim)
					if err != nil {
						fmt.Println("failed to decode: ", err)
						list.Error = err
						return
					}
					payloads <- cd
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
				decoded, err := DelimiterTag(elmt).Parse()
				if err != nil {
					return nil, err
				}
				delim = decoded
			}
		}
	}
}
