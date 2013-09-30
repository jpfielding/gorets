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
}

type MSystem struct {
	Version, Date string
	Id, Description string
	Comments string
}

type MResource struct {
	Resource, Version, Date string
	Id, Description string
	Comments string
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
	case "METADATA-CLASS":
	case "METADATA-TABLE":
	case "METADATA-LOOKUP":
	case "METADATA-LOOKUP_TYPE":
	}


	return &metadata, nil
}


func parseMSystem(response []byte) (*MSystem, error) {
	type XmlMSys struct {
		Version string `xml:"Version,attr"`
		Date string `xml:"Date,attr"`
		Comments string `xml:"COMMENTS"`
	}
	type XmlSys struct {
		SystemId string `xml:"SystemID,attr"`
		Description string `xml:"SystemDescription,attr"`
	}
	type XmlMSystem struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		MSystem XmlMSys `xml:"METADATA-SYSTEM"`
		System XmlSys `xml:"SYSTEM"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlMSystem{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	// transfer the contents to the public struct
	return &MSystem{
		Version: xms.MSystem.Version,
		Date: xms.MSystem.Date,
		Comments: xms.MSystem.Comments,
		Id: xms.System.SystemId,
		Description: xms.System.Description,
	}, nil
}
