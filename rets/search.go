package rets

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"context"
)

// Count is the type of serverside count response that should be included
const (
	// CountNone dont include a count
	CountNone = iota
	// CountIncluded include a count after the data
	CountIncluded = iota
	// CountOnly returns only the count
	CountOnly = iota
)

const (
	// StandardNames requests normalized naming
	StandardNames = iota
	// SystemNames (the default) requests whatever field names are used by the server
	SystemNames = iota
)

// SearchParams is the configuration for creating a SearchReqeust
type SearchParams struct {
	SearchType, // Property
	Class string // Residential

	HTTPFormEncodedValues bool // POST style http params

	Format, // 7.4.2 COMPACT | COMPACT-DECODED | STANDARD-XML | STANDARD-XML:dtd-version
	Select string

	// Payload should not be used with the format,select pair
	Payload string //The Client may request a specific XML format for the return set.

	// Query should be in the format specified by QueryType
	Query,
	QueryType string // DMQL2 is the standard option

	// RestrictedIndicator is the symbol to be used for fields that are blanked serverside (e.g. ####)
	RestrictedIndicator string

	StandardNames int // (0|1|)
	Count         int // (|0|1|2)
	Limit         int // <0 => "NONE"
	// Offset is properly started at 1, or left blank and assumed to have been 1
	Offset int
}

// SearchRequest holds the information needed to send a RETS request
type SearchRequest struct {
	URL,
	HTTPMethod string

	SearchParams
}

// PrepSearchRequest creates an http.Request from a SearchRequest
func PrepSearchRequest(r SearchRequest) (*http.Request, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	values := url.Query()
	// required
	values.Add("Class", r.Class)
	values.Add("SearchType", r.SearchType)

	// optional
	optionalString := OptionalStringValue(values)
	optionalString("Format", r.Format)
	optionalString("Select", r.Select)
	optionalString("Payload", r.Payload)
	optionalString("Query", r.Query)
	optionalString("QueryType", r.QueryType)
	optionalString("RestrictedIndicator", r.RestrictedIndicator)

	optionalInt := OptionalIntValue(values)
	if r.Count > 0 {
		optionalInt("Count", r.Count)
	}
	if r.Offset > 0 {
		optionalInt("Offset", r.Offset)
	}
	if r.StandardNames > 0 {
		optionalInt("StandardNames", r.StandardNames)
	}
	// limit is unique in that it can send a value of "NONE"
	switch {
	case r.Limit > 0:
		optionalInt("Limit", r.Limit)
	case r.Limit < 0:
		values.Add("Limit", "NONE")
	}

	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}

	// http POST style params
	if r.HTTPFormEncodedValues {
		req, err := http.NewRequest(method, url.String(), bytes.NewBufferString(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, err
	}
	// the standard query string style params here
	url.RawQuery = values.Encode()
	return http.NewRequest(method, url.String(), nil)
}

// SearchStream returns the raw stream from the RETS server response
func SearchStream(ctx context.Context, requester Requester, r SearchRequest) (io.ReadCloser, error) {
	req, err := PrepSearchRequest(r)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	return DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType)), nil
}

// countTag wraps an xml.StartElement to extract response Count information
type countTag xml.StartElement

// Parse ...
func (ct countTag) Parse() (int, error) {
	code, err := strconv.ParseInt(ct.Attr[0].Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(code), nil
}
