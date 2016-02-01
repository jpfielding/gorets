package client

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
)

const (
	// CountNone dont include a count
	CountNone = 0
	// CountAfter include a count after the data
	CountAfter = 1
	// CountOnly returns only the count
	CountOnly = 2
)

// TODO include standard names constants here

// CompactSearchResult ...
type CompactSearchResult struct {
	RetsResponse RetsResponse
	Count        int
	Delimiter    string
	Columns      []string
	MaxRows      bool

	Data   <-chan []string
	Errors chan error
}

// SearchRequest ...
type SearchRequest struct {
	URL,
	Class,
	SearchType,
	Format,
	Select,
	Query,
	QueryType,
	RestrictedIndicator,
	Payload,
	HTTPMethod string

	Count,
	// TODO NONE is a valid option, this needs to be modified
	Limit,
	Offset int
	BufferSize int
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

// SearchCompact if you set the wrong request Format you will get nothing back
func SearchCompact(requester Requester, ctx context.Context, r SearchRequest) (*CompactSearchResult, error) {
	body, err := SearchStream(requester, ctx, r)
	if err != nil {
		return nil, err
	}
	return NewCompactSearchResult(ctx, body, r.BufferSize)
}

// NewCompactSearchResult ...
func NewCompactSearchResult(ctx context.Context, body io.ReadCloser, bufferSize int) (*CompactSearchResult, error) {
	data := make(chan []string, bufferSize)
	errs := make(chan error)
	rets := RetsResponse{}
	result := &CompactSearchResult{
		Data:         data,
		Errors:       errs,
		RetsResponse: rets,
		MaxRows:      false,
	}

	parser := DefaultXMLDecoder(body, false)
	var buf bytes.Buffer

	// backgroundable processing of the data into our buffer
	bgDataProcessing := func() {
		// intentionally not closing errs
		defer close(data)
		defer body.Close()
		for {
			token, err := parser.Token()
			if err != nil {
				result.Errors <- err
				continue
			}
			switch t := token.(type) {
			case xml.StartElement:
				switch t.Name.Local {
				case "MAXROWS":
					result.MaxRows = true
				}
				// clear any accumulated data
				buf.Reset()
			case xml.EndElement:
				switch t.Name.Local {
				case "DATA":
					// need to select on both chans to avoid deadlock
					select {
					case <-ctx.Done():
						// no need to pipe the ctx err back as the caller already has it
						return
					case data <- ParseCompactRow(buf.String(), result.Delimiter):
					}
				case "RETS":
					return
				}
			case xml.CharData:
				bytes := xml.CharData(t)
				buf.Write(bytes)
			}
		}
	}

	// extract the basic content before delving into the data
	for {
		token, err := parser.Token()
		switch err {
		case nil:
			// nothing
		case io.EOF:
			return result, nil
		default:
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			buf.Reset()
			switch t.Name.Local {
			case "RETS", "RETS-STATUS":
				rets, err := ParseRetsResponseTag(t)
				if err != nil {
					return nil, err
				}
				result.RetsResponse = *rets
			case "COUNT":
				result.Count, err = ParseCountTag(t)
				if err != nil {
					return nil, err
				}
			case "DELIMITER":
				result.Delimiter, err = ParseDelimiterTag(t)
				if err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			switch name {
			case "COLUMNS":
				result.Columns = ParseCompactRow(buf.String(), result.Delimiter)
				go bgDataProcessing()
				return result, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			buf.Write(bytes)
		}
	}
}
