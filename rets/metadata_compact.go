package rets

import (
	"context"
	"encoding/xml"
	"io"
	"strings"
)

// CompactMetadataDelim is the only delimiter option for compact metadata
const CompactMetadataDelim = "	"

// CompactMetadata ...
type CompactMetadata struct {
	Rets     RetsResponse
	System   MSystem
	Elements map[string][]CompactData
}

// MSystem ...
type MSystem struct {
	Date, Version   string
	ID, Description string
	Comments        string
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
				cd, err := NewCompactData(t, parser, CompactMetadataDelim)
				if err != nil {
					return nil, err
				}
				metadata.Elements[cd.ID] = append(metadata.Elements[cd.ID], cd)
			}
		}
	}
}
