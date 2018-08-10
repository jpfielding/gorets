package rets

import (
	"context"
	"encoding/xml"
	"io"
)

// GetCompactMetadata ...
func GetCompactMetadata(ctx context.Context, requester Requester, r MetadataRequest) (*CompactMetadata, error) {
	r.Format = "COMPACT"
	body, err := MetadataStream(MetadataResponse(ctx, requester, r))
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
			// return io.EOF since we should have escaped cleanly with a RETS end tag
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				rets, err := ResponseTag(t).Parse()
				if err != nil {
					return nil, err
				}
				metadata.Response = *rets
			case "METADATA-SYSTEM":
				err := parser.DecodeElement(&metadata.MSystem, &t)
				if err != nil {
					return nil, err
				}
			default:
				cd, err := NewCompactData(t, parser, CompactMetadataDelim)
				if err != nil {
					return nil, err
				}
				metadata.Elements[cd.Element] = append(metadata.Elements[cd.Element], cd)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				return &metadata, nil

			}
		}
	}
}

// CompactMetadataDelim is the only delimiter option for compact metadata
const CompactMetadataDelim = "\t"

// CompactMetadata ...
type CompactMetadata struct {
	Response Response
	MSystem  CompactMSystem
	Elements map[string][]CompactData
}

// CompactMSystem ...
type CompactMSystem struct {
	Version  string        `xml:"Version,attr"`
	Date     string        `xml:"Date,attr"`
	Comments string        `xml:"COMMENTS"`
	System   CompactSystem `xml:"SYSTEM"`
}

// CompactSystem ...
type CompactSystem struct {
	ID             string `xml:"SystemID,attr"`
	Description    string `xml:"SystemDescription,attr"`
	TimeZoneOffset string `xml:",attr,omitempty"`
	MetadataID     string `xml:",attr,omitempty"`
}
