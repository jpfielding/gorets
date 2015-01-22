package gorets_metadata

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type MSystem struct {
	Version string
	Date    string
	System  System
}

type System struct {
	Id, Description string
	TimezoneOffset  string
	MetadataID      string
	// under the 'system'
	Comments          string
	ResourceVersion   string
	ResourceDate      string
	ForeignKeyVersion string
	ForeignKeyDate    string
	FilterVersion     string
	FilterDate        string
}

func (m *MSystem) InitFromCompact(p *xml.Decoder, t xml.StartElement) error {
	type Xsystem struct {
		SystemId          string `xml:"SystemID,attr"`
		Description       string `xml:"SystemDescription,attr"`
		TimeZoneOffset    string `xml:"TimeZoneOffset,attr"`
		MetadataID        string `xml:"MetadataID,attr"`
		ResourceVersion   string `xml:"ResourceVersion"`
		ResourceDate      string `xml:"ResourceDate"`
		ForeignKeyVersion string `xml:"ForeignKeyVersion"`
		ForeignKeyDate    string `xml:"ForeignKeyDate"`
		FilterVersion     string `xml:"FilterVersion"`
		FilterDate        string `xml:"FilterDate"`
	}
	type Xmsystem struct {
		Version  string  `xml:"Version,attr"`
		Date     string  `xml:"Date,attr"`
		System   Xsystem `xml:"SYSTEM"`
		Comments string  `xml:"COMMENTS"`
	}
	xms := Xmsystem{}
	err := p.DecodeElement(&xms, &t)
	if err != nil {
		return err
	}
	m.Version = xms.Version
	m.Date = xms.Date
	fmt.Println("COMMENTS: " + xms.Comments)
	m.System.Comments = strings.TrimSpace(xms.Comments)
	m.System.Id = xms.System.SystemId
	m.System.Description = xms.System.Description
	m.System.TimezoneOffset = xms.System.TimeZoneOffset
	m.System.MetadataID = xms.System.MetadataID
	m.System.ResourceVersion = xms.System.ResourceVersion
	m.System.ResourceDate = xms.System.ResourceDate
	m.System.ForeignKeyVersion = xms.System.ForeignKeyVersion
	m.System.ForeignKeyDate = xms.System.ForeignKeyDate
	m.System.FilterVersion = xms.System.FilterVersion
	m.System.FilterDate = xms.System.FilterDate
	return nil
}

func (m *MSystem) InitFromXml(p *xml.Decoder, t xml.StartElement) error {
	type system struct {
		SystemId          string `xml:"SystemID"`
		Description       string `xml:"SystemDescription"`
		TimeZoneOffset    string `xml:"TimeZoneOffset"`
		MetadataID        string `xml:"MetadataID"`
		Comments          string `xml:"Comments"`
		ResourceVersion   string `xml:"ResourceVersion"`
		ResourceDate      string `xml:"ResourceDate"`
		ForeignKeyVersion string `xml:"ForeignKeyVersion"`
		ForeignKeyDate    string `xml:"ForeignKeyDate"`
		FilterVersion     string `xml:"FilterVersion"`
		FilterDate        string `xml:"FilterDate"`
	}
	type msystem struct {
		Version string `xml:"Version,attr"`
		Date    string `xml:"Date,attr"`
		System  system `xml:"System"`
	}
	xms := msystem{}
	err := p.DecodeElement(&xms, &t)
	if err != nil {
		return err
	}
	m.System = System{}
	m.Version = xms.Version
	m.Date = xms.Date
	m.System.Comments = strings.TrimSpace(xms.System.Comments)
	m.System.Id = xms.System.SystemId
	m.System.Description = xms.System.Description
	m.System.TimezoneOffset = xms.System.TimeZoneOffset
	m.System.MetadataID = xms.System.MetadataID
	m.System.ResourceVersion = xms.System.ResourceVersion
	m.System.ResourceDate = xms.System.ResourceDate
	m.System.ForeignKeyVersion = xms.System.ForeignKeyVersion
	m.System.ForeignKeyDate = xms.System.ForeignKeyDate
	m.System.FilterVersion = xms.System.FilterVersion
	m.System.FilterDate = xms.System.FilterDate
	return nil
}
