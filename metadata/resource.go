package metadata

import (
	"encoding/xml"
)

type MResource struct {
	Version  string `xml:"Version,attr"`
	Date     string `xml:"Date,attr"`
	Resource []Resource
}

type Resource struct {
	ResourceID                  string
	StandardName                string
	VisibleName                 string
	Description                 string
	KeyField                    string
	ClassCount                  int
	ClassVersion                string
	ClassDate                   string
	ObjectVersion               string
	ObjectDate                  string
	SearchHelpVersion           string
	SearchHelpDate              string
	EditMaskVersion             string
	EditMaskDate                string
	LookupVersion               string
	LookupDate                  string
	UpdateHelpVersion           string
	UpdateHelpDate              string
	ValidationExpressionVersion string
	ValidationExpressionDate    string
	ValidationLookupVersion     string
	ValidationLookupDate        string
	ValidationExternalVersion   string
	ValidationExternalDate      string
}

func (m *MResource) InitFromXml(p *xml.Decoder, t xml.StartElement) error {
	err := p.DecodeElement(m, &t)
	if err != nil {
		return err
	}
	return nil
}

func (m *MResource) InitFromCompact(p *xml.Decoder, t xml.StartElement) error {
	cd, err := CompactData{}.Parse(p, t, "	")
	if err != nil {
		return err
	}
	m.Version = cd.Attrs["Version"]
	m.Date = cd.Attrs["Date"]
	for _, r := range cd.Data {
		res := Resource{}
		ApplyMap(r, &res)
		m.Resource = append(m.Resource, res)
	}
	return nil
}
