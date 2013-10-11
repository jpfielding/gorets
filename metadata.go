/**
	extraction of the data pieces describing a RETS system
 */
package gorets

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"strconv"
)

type Metadata struct {
	Rets RetsResponse
	System MSystem
	Resources MData
	Classes map[string]MData
	Tables map[string]MData
	Lookups map[string]MData
	LookupTypes map[string]MData
}

type MSystem struct {
	Date, Version string
	Id, Description string
	Comments string
}

/* the common structure */
type MData struct {
	Id, Date, Version string
	Columns []string
	Rows [][]string
}

/** cached lookup */
type Indexer func(col string, row int) string
/** create the cache */
func (m *MData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range m.Columns {
		index[c] = i
	}
	return func(col string, row int) string {
		return m.Rows[row][index[col]]
	}
}

type MetadataRequest struct {
	/* RETS request options */
	Url, Format, MType, Id string
}

func (s *Session) GetMetadata(r MetadataRequest) (*Metadata, error) {
	// required
	values := url.Values{}
	values.Add("Format", r.Format)
	values.Add("Type", r.MType)
	values.Add("ID", r.Id)

	req, err := http.NewRequest(s.HttpMethod, r.Url+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch r.Format {
	case "COMPACT":
		return parseMetadataCompactResult(resp.Body)
	case "STANDARD-XML":
		return parseMetadataStandardXml(resp.Body)
	}

	return nil, errors.New("unknows metadata format")
}

func parseMetadataCompactResult(body io.ReadCloser) (*Metadata,error) {
	parser := xml.NewDecoder(body)

	metadata := Metadata{}
	metadata.Classes = make(map[string]MData)
	metadata.Tables = make(map[string]MData)
	metadata.Lookups = make(map[string]MData)
	metadata.LookupTypes = make(map[string]MData)
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
				metadata.Rets.ReplyCode = int(code)
				metadata.Rets.ReplyText = attrs["replytext"]
			case "METADATA-SYSTEM":
				type XmlSystem struct {
					SystemId string `xml:"SystemID,attr"`
					Description string `xml:"SystemDescription,attr"`
				}
				type XmlMetadataSystem struct {
					Version string `xml:"Version,attr"`
					Date string `xml:"Date,attr"`
					System XmlSystem `xml:"SYSTEM"`
					Comments string `xml:"COMMENTS"`
				}
				xms := XmlMetadataSystem{}
				err := parser.DecodeElement(&xms, &t)
				if err != nil {
					return nil, err
				}
				metadata.System.Version = xms.Version
				metadata.System.Date = xms.Date
				metadata.System.Comments = strings.TrimSpace(xms.Comments)
				metadata.System.Id = xms.System.SystemId
				metadata.System.Description = xms.System.Description
			case "METADATA-RESOURCE", "METADATA-CLASS", "METADATA-TABLE", "METADATA-LOOKUP","METADATA-LOOKUP_TYPE":
				type XmlMetadataElement struct {
					Resource string `xml:"Resource,attr"`
					/* only valid for table */
					Class string `xml:"Class,attr"`
					/* only valid for lookup_type */
					Lookup string `xml:"Lookup,attr"`
					Version string `xml:"Version,attr"`
					Date string `xml:"Date,attr"`
					Columns string `xml:"COLUMNS"`
					Data []string `xml:"DATA"`
				}
				xme := XmlMetadataElement{}
				err := parser.DecodeElement(&xme,&t)
				if err != nil {
					fmt.Println("failed to decode: ", err)
					return nil, err
				}
				data := *extractMap(xme.Columns, xme.Data)
				data.Date = xme.Date
				if err != nil {
					return nil, err
				}
				data.Version = xme.Version
				data.Id = xme.Resource
				switch elmt.Name.Local {
				case "METADATA-RESOURCE":
					metadata.Resources = data
				case "METADATA-CLASS":
					metadata.Classes[data.Id] = data
				case "METADATA-TABLE":
					data.Id = xme.Resource +":"+ xme.Class
					metadata.Tables[data.Id] = data
				case "METADATA-LOOKUP":
					metadata.Lookups[data.Id] = data
				case "METADATA-LOOKUP_TYPE":
					data.Id = xme.Resource +":"+ xme.Lookup
					metadata.LookupTypes[data.Id] = data
				}
			}
		}
	}
	return &metadata, nil
}

const delim = "	"
/** extract a map of fields from columns and rows */
func extractMap(cols string, rows []string) (*MData) {
	data := MData{}
	// remove the first and last chars
	data.Columns = strings.Split(strings.Trim(cols,delim),delim)
	data.Rows = make([][]string, len(rows))
	// create each
	for i,line := range rows {
		data.Rows[i] = strings.Split(strings.Trim(line,delim),delim)
	}
	return &data
}


func parseMetadataStandardXml(body io.ReadCloser) (*Metadata,error) {
	return nil, errors.New("unsupported metadata format option")
}
