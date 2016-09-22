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

	retsStart, retsEnd := startend("RETS", map[string]string{
		"ReplyCode": fmt.Sprintf("%d", cm.Response.Code),
		"ReplyText": cm.Response.Text,
	})
	enc.EncodeToken(retsStart)
	msysStart, msysEnd := startend("METADATA-SYSTEM", map[string]string{
		"Date":    cm.MSystem.Date,
		"Version": cm.MSystem.Version,
	})
	enc.EncodeToken(msysStart)
	sysStart, sysEnd := startend("SYSTEM", map[string]string{
		"SystemID":          cm.MSystem.System.SystemID,
		"SystemDescription": cm.MSystem.System.Description,
	})
	enc.EncodeToken(sysStart)
	commentsStart, commentsEnd := startend("COMMENTS", map[string]string{})
	enc.EncodeToken(commentsStart)
	enc.EncodeToken(commentsEnd)

	// INNNER

	enc.EncodeToken(sysEnd)
	enc.EncodeToken(msysEnd)
	enc.EncodeToken(retsEnd)
	return enc.Flush()
}

func startend(name string, attrs map[string]string) (xml.StartElement, xml.EndElement) {
	elem := xml.Name{Local: name}
	var attr []xml.Attr
	for k, v := range attrs {
		attr = append(attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	return xml.StartElement{Name: elem, Attr: attr}, xml.EndElement{elem}
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
