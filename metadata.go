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
	MClasses MClasses
	MTables MTables
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

	// TOOD remove the needless repetition
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
		tmp, err := parseMClasses(body)
		if err != nil {
			return nil, err
		}
		metadata.MClasses = *tmp
	case "METADATA-TABLE":
		tmp, err := parseMTables(body)
		if err != nil {
			return nil, err
		}
		metadata.MTables = *tmp
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


const delim = "	"

/** extract a map of fields from columns and rows */
func extractMap(cols string, rows []string) ([]map[string]string) {
	// remove the first and last chars
	headers := strings.Split(strings.Trim(cols,delim),delim)
	fields := make([]map[string]string, len(rows))
	// create each
	for i,line := range rows {
		row := strings.Split(strings.Trim(line,delim),delim)
		fields[i] = make(map[string]string)
		for j, val := range row {
			fields[i][headers[j]] = val
		}
	}
	return fields
}

type MResources struct {
	Version, Date string
	Fields []map[string]string
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
		Info XmlResource `xml:"METADATA-RESOURCE"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlData{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	// remove the first and last chars
	rows := extractMap(xms.Info.Columns, xms.Info.Data)

	// transfer the contents to the public struct
	return &MResources{
		Version: xms.Info.Version,
		Date: xms.Info.Date,
		Fields: rows,
	}, nil
}

type MClasses struct {
	Version, Date string
	Fields []map[string]string
}

func parseMClasses(response []byte) (*MClasses, error) {
	type XmlClass struct {
		Resource string `xml:"Resource,attr"`
		Version string `xml:"Version,attr"`
		Date string `xml:"Date,attr"`
		Columns string `xml:"COLUMNS"`
		Data []string `xml:"DATA"`
	}
	type XmlData struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		Info XmlClass `xml:"METADATA-CLASS"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlData{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	// remove the first and last chars
	rows := extractMap(xms.Info.Columns, xms.Info.Data)

	// transfer the contents to the public struct
	return &MClasses{
		Version: xms.Info.Version,
		Date: xms.Info.Date,
		Fields: rows,
	}, nil
}

type MTables struct {
	Version, Date string
	Fields []map[string]string
}

func parseMTables(response []byte) (*MTables, error) {
	type XmlTable struct {
		Resource string `xml:"Resource,attr"`
		Version string `xml:"Version,attr"`
		Date string `xml:"Date,attr"`
		Columns string `xml:"COLUMNS"`
		Data []string `xml:"DATA"`
	}
	type XmlData struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		Info XmlTable `xml:"METADATA-TABLE"`
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false

	xms := XmlData{}
	err := decoder.Decode(&xms)
	if err != nil {
		return nil, err
	}

	// remove the first and last chars
	rows := extractMap(xms.Info.Columns, xms.Info.Data)

	// transfer the contents to the public struct
	return &MTables{
		Version: xms.Info.Version,
		Date: xms.Info.Date,
		Fields: rows,
	}, nil
}
