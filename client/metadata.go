/**
extraction of the data pieces describing a RETS system
*/
package client

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type Metadata struct {
	Rets        RetsResponse
	System      MSystem
	Resources   CompactData
	Classes     map[string]CompactData
	Tables      map[string]CompactData
	Lookups     map[string]CompactData
	LookupTypes map[string]CompactData
}

type MSystem struct {
	Date, Version   string
	Id, Description string
	Comments        string
}

type MetadataRequest struct {
	/* RETS request options */
	URL, HTTPMethod, Format, MType, Id string
}

func (s *Session) GetMetadata(ctx context.Context, r MetadataRequest) (*Metadata, error) {
	// required
	values := url.Values{}
	values.Add("Format", r.Format)
	values.Add("Type", r.MType)
	values.Add("ID", r.Id)

	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	// TODO use a URL object then properly append to it
	req, err := http.NewRequest(method, r.URL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, &s.Client, req)
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

func parseMetadataCompactResult(body io.ReadCloser) (*Metadata, error) {
	parser := xml.NewDecoder(body)

	metadata := Metadata{}
	metadata.Classes = make(map[string]CompactData)
	metadata.Tables = make(map[string]CompactData)
	metadata.Lookups = make(map[string]CompactData)
	metadata.LookupTypes = make(map[string]CompactData)
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
				type XmlSystem struct {
					SystemId    string `xml:"SystemID,attr"`
					Description string `xml:"SystemDescription,attr"`
				}
				type XmlMetadataSystem struct {
					Version  string    `xml:"Version,attr"`
					Date     string    `xml:"Date,attr"`
					System   XmlSystem `xml:"SYSTEM"`
					Comments string    `xml:"COMMENTS"`
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
			case "METADATA-RESOURCE", "METADATA-CLASS", "METADATA-TABLE", "METADATA-LOOKUP", "METADATA-LOOKUP_TYPE":
				data, err := ParseMetadataCompactDecoded(t, parser, "	")
				if err != nil {
					return nil, err
				}
				if data == nil {
					continue
				}
				switch t.Name.Local {
				case "METADATA-RESOURCE":
					metadata.Resources = *data
				case "METADATA-CLASS":
					metadata.Classes[data.Id] = *data
				case "METADATA-TABLE":
					metadata.Tables[data.Id] = *data
				case "METADATA-LOOKUP":
					metadata.Lookups[data.Id] = *data
				case "METADATA-LOOKUP_TYPE":
					metadata.LookupTypes[data.Id] = *data
				}
			}
		}
	}
}

func parseMetadataStandardXml(body io.ReadCloser) (*Metadata, error) {
	defer body.Close()
	ioutil.ReadAll(body)
	return nil, errors.New("unsupported metadata format option")
}
