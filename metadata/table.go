package metadata

import (
	"encoding/xml"
)

// MTable is the container for 'Table' (fields) metadata compoents
type MTable struct {
	Version  string  `xml:",attr"`
	Date     string  `xml:",attr"`
	Resource string  `xml:",attr"`
	Class    string  `xml:",attr"`
	Table    []Table `xml:"Field"`
}

// Table is the metadata descriptor of fields
type Table struct {
	MetadataEntryID string
	SystemName      string
	StandardName    string
	LongName        string
	DBName          string
	ShortName       string
	MaximumLength   string
	// Boolean, Character, Date, DateTime, Time, Tiny, Small, Int, Long, Decimal
	DataType   string
	Precision  string
	Searchable string
	// Number, Currency, Lookup, LookupMulti, URI
	Interpretation string
	// Left, Right, Center, Justify
	Alignment    string
	UseSeparator string
	EditMaskID   string
	LookupName   string
	MaxSelect    string
	// Feet | Meters | SqFt | SqMeters | Acres | Hectares
	Units              string
	Index              string
	Minimum            string
	Maximum            string
	Default            string
	Required           string
	SearchHelpID       string
	Unique             string
	ForeignKeyName     string
	ForeignKeyField    string
	InKeyIndex         string
	FilterParentField  string
	DefaultSearchOrder string
	// UPPER, LOWER, EXACT, MIXED
	Case string
}

// InitFromXml parses the xml form of this metdata component
func (m *MTable) InitFromXml(p *xml.Decoder, t xml.StartElement) error {
	err := p.DecodeElement(m, &t)
	if err != nil {
		return err
	}
	return nil
}

// InitFromCompact parses the compact form of this metdata component
func (m *MTable) InitFromCompact(p *xml.Decoder, t xml.StartElement) error {
	cd, err := CompactData{}.Parse(p, t, "	")
	if err != nil {
		return err
	}
	m.Version = cd.Attrs["Version"]
	m.Date = cd.Attrs["Date"]
	m.Resource = cd.Attrs["Resource"]
	m.Class = cd.Attrs["Class"]
	for _, r := range cd.Data {
		res := Table{}
		ApplyMap(r, &res)
		m.Table = append(m.Table, res)
	}
	return nil
}
