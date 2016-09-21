package rets

import (
	"context"
	"encoding/xml"
	"io"
)

// CompactMetadataDelim is the only delimiter option for compact metadata
const CompactMetadataDelim = "	"

// CompactMetadata ...
type CompactMetadata struct {
	Rets     RetsResponse
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
	SystemID    string `xml:"SystemID,attr"`
	Description string `xml:"SystemDescription,attr"`
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
				err := parser.DecodeElement(&metadata.MSystem, &t)
				if err != nil {
					return nil, err
				}
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
