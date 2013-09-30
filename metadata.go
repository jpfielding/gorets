/**
	extraction of the data pieces describing a RETS system
 */
package gorets

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"encoding/xml"
)

type MSystem struct {
	Version, Date string
	Id, Description string
	Comments string
}

type Metadata struct {
	MSystem MSystem
}

func (s *Session) GetMetadata(url, format, id, mtype string) (*Metadata, error) {
	qs := fmt.Sprintf("Format=%s",format)
	if id != "" {
		qs = qs +"&"+ fmt.Sprintf("ID=%s",id)
	}
	if mtype != "" {
		qs = qs +"&"+ fmt.Sprintf("Type=%s",mtype)
	}

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

	switch mtype {
	case "METADATA-SYSTEM":
		return parseMSystem(body)
	case "METADATA-RESOURCE":
	case "METADATA-CLASS":
	case "METADATA-TABLE":
	case "METADATA-LOOKUP":
	case "METADATA-LOOKUP_TYPE":
	}
	metadata := &Metadata{}
	return metadata, nil
}


func parseMSystem(response []byte) (*Metadata, error) {
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
	if err != nil && err != io.EOF {
		return nil, err
	}

	// transfer the contents to the public struct
	mSystem := MSystem{
		Version: xms.MSystem.Version,
		Date: xms.MSystem.Date,
		Comments: xms.MSystem.Comments,
		Id: xms.System.SystemId,
		Description: xms.System.Description,
	}

	return &Metadata{
		MSystem: mSystem,
	}, nil
}
