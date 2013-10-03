/**
	provides the searching core

	see minidom style processing here:
		http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
 */
package gorets

import (
	"encoding/xml"
	"encoding/hex"
	"errors"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)


/* counts */
const (
	COUNT_NONE = 0
	COUNT_AFTER = 1
	COUNT_ONLY = 2
)
/* field naming */
const (
	STANDARD_NAMES_OFF = 0
	STANDARD_NAMES_ON = 0
)

type SearchResult struct {
	RetsResponse RetsResponse
	Count int64
	Delimiter string
	Columns []string
	Data <-chan []string
	MaxRows bool
	ProcessingFailure error
}

func (m *SearchResult) Index() (map[string]int) {
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
	return func(row []string) (filtered []string){
		tmp := make([]string, len(cols))
		for i,c := range cols {
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
	Limit,
	Offset int
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
func (s *Session) Search(r SearchRequest) (*SearchResult, error) {
	// required
	values := url.Values{}
	values.Add("Class", r.Class)
	values.Add("SearchType", r.SearchType)

	// optional
	optionalString := func (name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}
	optionalString("Format", r.Format)
	optionalString("Select", r.Select)
	optionalString("Query", r.Query)
	optionalString("QueryType", r.QueryType)
	optionalString("RestrictedIndicator", r.RestrictedIndicator)
	optionalString("Payload", r.Payload)

	optionalInt := func (name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d",value))
		}
	}
	optionalInt("Count", r.Count)
	optionalInt("Limit", r.Limit)
	optionalInt("Offset", r.Offset)

	req, err := http.NewRequest(s.HttpMethod, fmt.Sprintf("%s?%s",r.Url,values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	switch r.Format {
	case "COMPACT-DECODED", "COMPACT":
		return parseCompactResult(resp.Body)
	case "STANDARD-XML":
		panic("not yet supported!")
	}
	return nil, nil
}


func parseCompactResult(body io.ReadCloser) (*SearchResult,error) {
	// TODO parameterize this buffer size
	data := make(chan []string,100)
	rets := RetsResponse{}
	result := SearchResult{
		Data: data,
		RetsResponse: rets,
		MaxRows: false,
	}

	parser := xml.NewDecoder(body)
	var buf bytes.Buffer

	// backgroundable processing of the data into our buffer
	dataProcessing := func() {
		delim := result.Delimiter
		// this channel needs to be closed or the caller can infinite loop
		defer close(data)
		// this is the web socket that needs addressed
		defer body.Close()
		// extract the data
		for {
			// TODO figure out a kill switch for this
			token, err := parser.Token()
			if err != nil {
				result.ProcessingFailure = err
				break
			}
			switch t := token.(type) {
			case xml.StartElement:
				elmt := xml.StartElement(t)
				name := elmt.Name.Local
				switch name {
				case "MAXROWS":
					result.MaxRows = true
				}
				// clear any accumulated data
				buf.Reset()
			case xml.EndElement:
				elmt := xml.EndElement(t)
				name := elmt.Name.Local
				switch name {
				case "DATA":
					data <- strings.Split(strings.Trim(buf.String(), delim), delim)
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
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			switch name {
			case "RETS":
				attrs := make(map[string]string)
				for _,v := range elmt.Attr {
					attrs[strings.ToLower(v.Name.Local)] = v.Value
				}
				// TODO dont rely on position, search by name
				code,err := strconv.ParseInt(attrs["replycode"],10,16)
				if err != nil {
					return nil, err
				}
				result.RetsResponse.ReplyCode = int(code)
				result.RetsResponse.ReplyText = attrs["replytext"]
			case "COUNT":
				result.Count,err = strconv.ParseInt(elmt.Attr[0].Value, 10, 64)
				if err != nil {
					return nil, err
				}
			case "DELIMITER":
				decoded,err := hex.DecodeString(elmt.Attr[0].Value)
				if err != nil {
					return nil, err
				}
				result.Delimiter = string(decoded)
			}
		case xml.EndElement:
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			switch name {
			case "COLUMNS":
				result.Columns = strings.Split(strings.Trim(buf.String(), result.Delimiter), result.Delimiter)
				go dataProcessing()
				return &result, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			buf.Write(bytes)
		}
	}

	return nil, errors.New("failed to parse rets response")
}

func parseStandardXml(body *io.ReadCloser) (*SearchResult,error) {
	return nil, nil
}

