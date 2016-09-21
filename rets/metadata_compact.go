package rets

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
)

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
				rets, err := ParseResponse(t)
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
				metadata.Elements[cd.ID] = append(metadata.Elements[cd.ID], cd)
			}
		}
	}
}

// CompactMetadataDelim is the only delimiter option for compact metadata
const CompactMetadataDelim = "	"

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
	SystemID    string `xml:"SystemID,attr"`
	Description string `xml:"SystemDescription,attr"`
}

// ToXML ...
func (cm CompactMetadata) ToXML(out io.Writer) error {
	// TODO do we need to include the xml.Header directly?
	out.Write([]byte(xml.Header))
	// os.Stdout
	enc := xml.NewEncoder(out)
	enc.Indent("  ", "    ")

	rets := xml.Name{Local: "RETS"}
	start := xml.StartElement{
		Name: rets,
		Attr: []xml.Attr{
			xml.Attr{Name: xml.Name{Local: "ReplyCode"}, Value: fmt.Sprintf("%d", cm.Response.Code)},
			xml.Attr{Name: xml.Name{Local: "ReplyText"}, Value: cm.Response.Text},
		},
	}
	enc.EncodeToken(start)
	defer enc.EncodeToken(xml.EndElement{rets})
	return nil
}

// MetadataLookup .. TODO figure out what to do with this
var MetadataLookup = map[string]string{
	"METADATA-SYSTEM":                     "SYSTEM",
	"METADATA-FOREIGN_KEY":                "ForeignKey",
	"METADATA-FILTER":                     "Filter",
	"METADATA-FILTER_TYPE":                "FilterType",
	"METADATA-RESOURCE":                   "Resource",
	"METADATA-CLASS":                      "Class",
	"METADATA-TABLE":                      "Field",
	"METADATA-UPDATE":                     "Update",
	"METADATA-UPDATE_TYPE":                "UpdateType",
	"METADATA-COLUMN_GROUP_SET":           "ColumnGroupSet",
	"METADATA-COLUMN_GROUP":               "ColumnGroup",
	"METADATA-COLUMN_GROUP_CONTROL":       "ColumnGroupControl",
	"METADATA-COLUMN_GROUP_TABLE":         "ColumnGroupTable",
	"METADATA-COLUMN_GROUP_NORMALIZATION": "ColumnGroupNormalization",

	"METADATA-OBJECT":                 "Object",
	"METADATA-SEARCHHELP":             "SearchHelp",
	"METADATA-EDITMASK":               "EditMask",
	"METADATA-UPDATEHELP":             "UpdateHelp",
	"METADATA-LOOKUP":                 "Lookup",
	"METADATA-LOOKUP_TYPE":            "LookupType",
	"METADATA-VALIDATION_LOOKUP":      "ValidationLookup",
	"METADATA-VALIDATION_LOOKUP_TYPE": "ValidationLookupType",

	"METADATA-VALIDATION_EXTERNAL":      "ValidationExternal",
	"METADATA-VALIDATION_EXTERNAL_TYPE": "ValidationExternalType",
	"METADATA-VALIDATION_EXPRESSION":    "ValidationExpression",
}
