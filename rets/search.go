package rets

import (
	"io"
	"net/http"
	"net/url"

	"context"
)

const (
	// CountNone dont include a count
	CountNone = 0
	// CountIncluded include a count after the data
	CountIncluded = 1
	// CountOnly returns only the count
	CountOnly = 2
)

// TODO include standard names constants here

// SearchParams ...
type SearchParams struct {
	Class,
	SearchType,
	Format, // 7.4.2 COMPACT | COMPACT-DECODED | STANDARD-XML
	Select,
	Query,
	QueryType,
	RestrictedIndicator,
	Payload string //The Client may request a specific XML format for the return set.

	Count,
	Limit, // <0 => "NONE"
	Offset int
}

// SearchRequest ...
type SearchRequest struct {
	URL,
	HTTPMethod string

	SearchParams

	BufferSize int // TODO unused atm
}

// SearchStream ...
func SearchStream(requester Requester, ctx context.Context, r SearchRequest) (io.ReadCloser, error) {
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
	optionalString("Query", r.Query)
	optionalString("QueryType", r.QueryType)
	optionalString("RestrictedIndicator", r.RestrictedIndicator)
	optionalString("Payload", r.Payload)

	optionalInt := OptionalIntValue(values)
	optionalInt("Count", r.Count)
	if r.Offset > 0 {
		optionalInt("Offset", r.Offset)
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

	url.RawQuery = values.Encode()

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	return DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType)), nil
}
