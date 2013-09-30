/**
	extraction of the data pieces describing a RETS system
 */
package gorets

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Metadata struct {
	MSystem MSystem
	MResources MResources
}

func (s *Session) GetMetadata(url, format, id, mtype string) (*Metadata, error) {
	qs := fmt.Sprintf("Format=%s",format)
	qs = qs +"&"+ fmt.Sprintf("Type=%s",mtype)
	qs = qs +"&"+ fmt.Sprintf("ID=%s",id)

	req, err := http.NewRequest(s.HttpMethod, url+"?"+qs, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	metadata := Metadata{}

	switch strings.ToUpper(mtype) {
	case "METADATA-SYSTEM":
		tmp, err := parseMSystem(body)
		if err != nil {
			return nil, err
		}
		metadata.MSystem = *tmp
	case "METADATA-RESOURCE":
		tmp, err := parseMResources(body)
		if err != nil {
			return nil, err
		}
		metadata.MResources = *tmp
	case "METADATA-CLASS":
	case "METADATA-TABLE":
	case "METADATA-LOOKUP":
	case "METADATA-LOOKUP_TYPE":
	}


	return &metadata, nil
}

type MSystem struct {
	Version, Date string
	Id, Description string
	Comments string
}

func parseMSystem(response []byte) (*MSystem, error) {
	type XmlMSystem struct {
		Version string `xml:"Version,attr"`
		Date string `xml:"Date,attr"`
		Comments string `xml:"COMMENTS"`
	}
	type XmlSystem struct {
		SystemId string `xml:"SystemID,attr"`
		Description string `xml:"SystemDescription,attr"`
	}
	type XmlData struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		MSystem XmlMSystem `xml:"METADATA-SYSTEM"`
		System XmlSystem `xml:"SYSTEM"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlData{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	// transfer the contents to the public struct
	return &MSystem{
		Version: xms.MSystem.Version,
		Date: xms.MSystem.Date,
		Comments: strings.TrimSpace(xms.MSystem.Comments),
		Id: xms.System.SystemId,
		Description: xms.System.Description,
	}, nil
}

type MResource struct {
	Version, Date string
	Fields map[string]string
}

type MResources struct {
	Version, Date string
	MResources []MResource
}

func parseMResources(response []byte) (*MResources, error) {
	type XmlResource struct {
		Version string `xml:"Version,attr"`
		Date string `xml:"Date,attr"`
		Columns string `xml:"COLUMNS"`
		Data []string `xml:"DATA"`
	}
	type XmlData struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		ResourceInfo XmlResource `xml:"METADATA-RESOURCE"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlData{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	tab := "	"
	// remove the first and last chars
	headers := strings.Split(strings.Trim(xms.ResourceInfo.Columns,tab),tab)
	resources := make([]MResource, len(xms.ResourceInfo.Data))
	// create each
	for i,line := range xms.ResourceInfo.Data {
		row := strings.Split(strings.Trim(line,tab),tab)
		resources[i].Fields = make(map[string]string)
		for j, val := range row {
			resources[i].Fields[headers[j]] = val
		}
	}

	// transfer the contents to the public struct
	return &MResources{
		Version: xms.ResourceInfo.Version,
		Date: xms.ResourceInfo.Date,
		MResources: resources,
	}, nil
}

