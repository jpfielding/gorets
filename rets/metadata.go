package rets

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"context"
)

// CompactMetadata ...
type CompactMetadata struct {
	Rets     RetsResponse
	System   MSystem
	Elements map[string][]CompactData
}

// CompactData is the common compact decoded structure
type CompactData struct {
	ID, Date, Version string
	Columns           Row
	Rows              []Row
}

// Indexer provices cached lookup for CompactData
type Indexer func(col string, row int) string

// Indexer create the cache
func (cd *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range cd.Columns {
		index[c] = i
	}
	return func(col string, row int) string {
		return cd.Rows[row][index[col]]
	}
}

func (cm CompactMetadata) find(name string) map[string]CompactData {
	classes := make(map[string]CompactData)
	for _, cd := range cm.Elements[name] {
		classes[cd.ID] = cd
	}
	return classes
}

// MSystem ...
type MSystem struct {
	Date, Version   string
	ID, Description string
	Comments        string
}

// MetadataRequest ...
type MetadataRequest struct {
	// RETS request options
	URL, HTTPMethod, Format, MType, ID string
}

// MetadataStream ...
func MetadataStream(requester Requester, ctx context.Context, r MetadataRequest) (io.ReadCloser, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	values := url.Query()
	// required
	values.Add("Format", r.Format)
	values.Add("Type", r.MType)
	values.Add("ID", r.ID)

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

// GetCompactMetadata ...
func GetCompactMetadata(requester Requester, ctx context.Context, r MetadataRequest) (*CompactMetadata, error) {
	r.Format = "COMPACT"
	body, err := MetadataStream(requester, ctx, r)
	if err != nil {
		return nil, err
	}
	return ParseMetadataCompactResult(body)
}

// ParseMetadataCompactResult ...
func ParseMetadataCompactResult(body io.ReadCloser) (*CompactMetadata, error) {
	defer body.Close()
	parser := DefaultXMLDecoder(body, false)

	metadata := CompactMetadata{}
	metadata.Elements = make(map[string][]CompactData)
	for {
		token, err := parser.Token()
		if err != nil {
			if err == io.EOF {
				return &metadata, nil
			}
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "RETS", "RETS-STATUS":
				rets, err := ParseRetsResponseTag(t)
				if err != nil {
					return nil, err
				}
				metadata.Rets = *rets
			case "METADATA-SYSTEM":
				type xmlSystem struct {
					SystemID    string `xml:"SystemID,attr"`
					Description string `xml:"SystemDescription,attr"`
				}
				type xmlMetadataSystem struct {
					Version  string    `xml:"Version,attr"`
					Date     string    `xml:"Date,attr"`
					System   xmlSystem `xml:"SYSTEM"`
					Comments string    `xml:"COMMENTS"`
				}
				xms := xmlMetadataSystem{}
				err := parser.DecodeElement(&xms, &t)
				if err != nil {
					return nil, err
				}
				metadata.System.Version = xms.Version
				metadata.System.Date = xms.Date
				metadata.System.Comments = strings.TrimSpace(xms.Comments)
				metadata.System.ID = xms.System.SystemID
				metadata.System.Description = xms.System.Description
			default:
				data, err := ParseMetadataCompactDecoded(t, parser, "	")
				if err != nil {
					return nil, err
				}
				if data == nil {
					continue
				}
				metadata.Elements[t.Name.Local] = append(metadata.Elements[t.Name.Local], *data)
			}
		}
	}
}

// ParseMetadataCompactDecoded ...
func ParseMetadataCompactDecoded(start xml.StartElement, parser *xml.Decoder, delim string) (*CompactData, error) {
	// XmlMetadataElement is the simple extraction tool for our data
	type XMLMetadataElement struct {
		Resource string `xml:"Resource,attr"`
		/* only valid for table */
		Class string `xml:"Class,attr"`
		/* only valid for lookup_type */
		Lookup  string   `xml:"Lookup,attr"`
		Version string   `xml:"Version,attr"`
		Date    string   `xml:"Date,attr"`
		Columns string   `xml:"COLUMNS"`
		Data    []string `xml:"DATA"`
	}
	xme := XMLMetadataElement{}
	err := parser.DecodeElement(&xme, &start)
	if err != nil {
		fmt.Println("failed to decode: ", err)
		return nil, err
	}
	if xme.Columns == "" {
		return nil, nil
	}
	data := *extractMap(xme.Columns, xme.Data, delim)
	data.Date = xme.Date
	data.Version = xme.Version
	data.ID = xme.Resource
	if xme.Class != "" {
		data.ID = xme.Resource + ":" + xme.Class
	}
	if xme.Lookup != "" {
		data.ID = xme.Resource + ":" + xme.Lookup
	}

	return &data, nil
}

/** extract a map of fields from columns and rows */
func extractMap(cols string, rows []string, delim string) *CompactData {
	data := CompactData{}
	// remove the first and last chars
	data.Columns = CompactRow(cols).Parse(delim)
	data.Rows = make([]Row, len(rows))
	// create each
	for i, line := range rows {
		data.Rows[i] = CompactRow(line).Parse(delim)
	}
	return &data
}
