package rets

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"

	"context"
)

// PayloadListRequest ...
type PayloadListRequest struct {
	URL, HTTPMethod, ID string
}

// PrepGetPayloadList creates an http.Request from a PayloadListRequest
func PrepGetPayloadList(r PayloadListRequest) (*http.Request, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	values := url.Query()
	// required
	values.Add("ID", r.ID)

	method := DefaultHTTPMethod
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	url.RawQuery = values.Encode()

	return http.NewRequest(method, url.String(), nil)
}

// GetPayloadList ...
func GetPayloadList(requester Requester, ctx context.Context, r PayloadListRequest) (PayloadList, error) {
	req, err := PrepGetPayloadList(r)
	if err != nil {
		return PayloadList{}, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return PayloadList{}, err
	}
	defer resp.Body.Close()

	return NewPayloadList(resp.Body)
}

// PayloadList ...
type PayloadList struct {
	Response Response

	body   io.ReadCloser
	parser *xml.Decoder
	delim  string
}

// EachPayload ...
type EachPayload func(CompactData, error) error

// ForEach ...
func (pl PayloadList) ForEach(each EachPayload) error {
	defer pl.body.Close()

	for {
		token, err := pl.parser.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "RETSPayloadList":
				err := each(NewCompactData(t, pl.parser, pl.delim))
				if err != nil {
					return err
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				return nil
			}
		}
	}
}

// NewPayloadList parse a stream and reads PayloadLists
func NewPayloadList(body io.ReadCloser) (PayloadList, error) {
	parser := xml.NewDecoder(body)

	// return a composite result and offer walk/close options
	pl := PayloadList{
		body:   body,
		parser: parser,
	}
	for {
		token, err := parser.Token()
		if err != nil {
			return pl, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				rets, er := ResponseTag(t).Parse()
				if err != nil {
					return pl, er
				}
				pl.Response = *rets
				return pl, nil
			case "DELIMITER":
				pl.delim, err = DelimiterTag(t).Parse()
				if err != nil {
					return pl, err
				}
			}
		}
	}
}
