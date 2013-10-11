package gorets

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"strconv"
)

type Payload struct {
	Resource, Class, Version, Date string
	Columns []string
	Rows [][]string
}

type PayloadList struct {
	Rets RetsResponse
	Error error
	Payloads <-chan Payload
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
	payloads := make(chan Payload)
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
				elmt := xml.StartElement(t)
				switch elmt.Name.Local {
				case "RETSPayloadList":
					type XmlPayloadList struct {
						Resource string `xml:"Resource,attr"`
						/* only valid for table */
						Class string `xml:"Class,attr"`
						Version string `xml:"Version,attr"`
						Date string `xml:"Date,attr"`
						Columns string `xml:"COLUMNS"`
						Data []string `xml:"DATA"`
					}
					xme := XmlPayloadList{}
					err := parser.DecodeElement(&xme,&t)
					if err != nil {
						fmt.Println("failed to decode: ", err)
						list.Error = err
						return
					}
					payload := Payload{
						Resource: xme.Resource,
						Class: xme.Class,
						Date: xme.Date,
						Version: xme.Version,
						Columns: strings.Split(strings.Trim(xme.Columns, delim), delim),
					}
					payload.Rows = make([][]string, len(xme.Data))
					// create each
					for i,line := range xme.Data {
						payload.Rows[i] = strings.Split(strings.Trim(line,delim),delim)
					}
					payloads <- payload
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
			case "RETS":
				attrs := make(map[string]string)
				for _,v := range elmt.Attr {
					attrs[strings.ToLower(v.Name.Local)] = v.Value
				}
				code,err := strconv.ParseInt(attrs["replycode"],10,16)
				if err != nil {
					return nil, err
				}
				if code != 0 {
					return nil, errors.New(attrs["replytext"])
				}
				list.Rets = RetsResponse{
					ReplyCode:  int(code),
					ReplyText: attrs["replytext"],
				}
				go dataProcessing()
				return &list, nil
			case "DELIMITER":
				decoded,err := hex.DecodeString(elmt.Attr[0].Value)
				if err != nil {
					return nil, err
				}
				delim = string(decoded)
			}
		}
	}
	return nil,errors.New("could not find rets response")
}
