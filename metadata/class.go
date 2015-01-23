package metadata

import "encoding/xml"

type MClass struct {
	Version  string `xml:"Version,attr"`
	Date     string `xml:"Date,attr"`
	Resource string `xml:"Resource,attr"`
	Class    []Class
}

type Class struct {
	ClassName        string
	StandardName     string
	VisibleName      string
	Description      string
	TableVersion     string
	TableDate        string
	UpdateVersion    string
	UpdateDate       string
	ClassTimeStamp   string
	DeletedFlagField string
	DeletedFlagValue string
	HasKeyIndex      string // Bool?
	OffsetSupport    string // Bool?
}

func (m *MClass) InitFromXml(p *xml.Decoder, t xml.StartElement) error {
	err := p.DecodeElement(m, &t)
	if err != nil {
		return err
	}
	return nil
}

func (m *MClass) InitFromCompact(p *xml.Decoder, t xml.StartElement) error {
	cd, err := CompactData{}.Parse(p, t, "	")
	if err != nil {
		return err
	}
	m.Version = cd.Attrs["Version"]
	m.Date = cd.Attrs["Date"]
	m.Resource = cd.Attrs["Resource"]
	for _, r := range cd.Data {
		res := Class{}
		res.ClassName = r["ClassName"]
		res.StandardName = r["StandardName"]
		res.VisibleName = r["VisibleName"]
		res.Description = r["Description"]
		res.TableVersion = r["TableVersion"]
		res.TableDate = r["TableDate"]
		res.UpdateVersion = r["UpdateVersion"]
		res.UpdateDate = r["UpdateDate"]
		res.ClassTimeStamp = r["ClassTimeStamp"]
		res.DeletedFlagField = r["DeletedFlagField"]
		res.DeletedFlagValue = r["DeletedFlagValue"]
		res.HasKeyIndex = r["HasKeyIndex"]
		res.OffsetSupport = r["OffsetSupport"]

		m.Class = append(m.Class, res)
	}
	return nil
}
