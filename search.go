/**
provides the searching core

see minidom style processing here:
	http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
*/
package gorets_client

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

/* counts */
const (
	COUNT_NONE  = 0
	COUNT_AFTER = 1
	COUNT_ONLY  = 2
)

/* field naming */
const (
	STANDARD_NAMES_OFF = 0
	STANDARD_NAMES_ON  = 0
)

type SearchResult struct {
	RetsResponse      RetsResponse
	Count             int
	Delimiter         string
	Columns           []string
	Data              <-chan []string
	MaxRows           bool
	ProcessingFailure error
}

func (m *SearchResult) Index() map[string]int {
	index := make(map[string]int)
	for i, c := range m.Columns {
		index[c] = i
	}
	return index
}

/** cached filtering */
type ColumnFilter func(row []string) (filtered []string)

/** create the cache */
func (m *SearchResult) FilterTo(cols []string) ColumnFilter {
	index := m.Index()
	return func(row []string) (filtered []string) {
		tmp := make([]string, len(cols))
		for i, c := range cols {
			tmp[i] = row[index[c]]
		}
		return tmp
	}
}

type SearchRequest struct {
	Url,
	Class,
	SearchType,
	Format,
	Select,
	Query,
	QueryType,
	RestrictedIndicator,
	Payload string

	Count,
	// TODO NONE is a valid option, this needs to be modified
	Limit,
	Offset int
	BufferSize int
}

/*
	GET /platinum/search?
	Class=ALL&
	Count=1&
	Format=COMPACT-DECODED&
	Limit=10&
	Offset=50&
	Query=%28%28LocaleListingStatus%3D%7CACTIVE-CORE%2CCNTG%2FKO-CORE%2CCNTG%2FNO+KO-CORE%2CAPP+REG-CORE%29%2C%7E%28VOWList%3D0%29%29&
	QueryType=DMQL2&
	SearchType=Property
*/
func (s *Session) Search(quit <-chan struct{}, r SearchRequest) (*SearchResult, error) {
	// required
	values := url.Values{}
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

	req, err := http.NewRequest(s.HttpMethod, fmt.Sprintf("%s?%s", r.Url, values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	switch r.Format {
	case "COMPACT-DECODED", "COMPACT":
		return parseCompactResult(quit, resp.Body, r.BufferSize)
	case "STANDARD-XML":
		panic("not yet supported!")
	}
	return nil, nil
}

func parseCompactResult(quit <-chan struct{}, body io.ReadCloser, bufferSize int) (*SearchResult, error) {
	if bufferSize <= 0 {
		bufferSize = 1
	}
	data := make(chan []string, bufferSize)
	rets := RetsResponse{}
	result := SearchResult{
		Data:         data,
		RetsResponse: rets,
		MaxRows:      false,
	}

	parser := xml.NewDecoder(body)
	var buf bytes.Buffer

	// backgroundable processing of the data into our buffer
	bgDataProcessing := func() {
		defer close(data)
		defer body.Close()
		for {
			token, err := parser.Token()
			if err != nil {
				result.ProcessingFailure = err
				break
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
					select {
					case data <- ParseCompactRow(buf.String(), result.Delimiter):
					case <-quit:
						return
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
		if err != nil {
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
				return &result, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			buf.Write(bytes)
		}
	}

	return nil, errors.New("failed to parse rets response")
}

func parseStandardXml(body *io.ReadCloser) (*SearchResult, error) {
	return nil, nil
}
