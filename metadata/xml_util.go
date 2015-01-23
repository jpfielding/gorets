package metadata

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// advances the cursor to the named xml.StartElement
func AdvanceToStartElem(parser *xml.Decoder, start string) (xml.StartElement, error) {
	for {
		token, err := parser.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			switch t.Name.Local {
			case start:
				return t, nil
			}
		}
	}
}

type CompactData struct {
	Attrs map[string]string
	Data  []map[string]string
}

func (cd CompactData) Parse(d *xml.Decoder, s xml.StartElement, delim string) (*CompactData, error) {
	cd.Attrs = make(map[string]string)
	cd.Data = make([]map[string]string, 0)
	// extract the meta
	for _, a := range s.Attr {
		cd.Attrs[a.Name.Local] = a.Value
	}
	// extract the rows
	type XmlMetadataElement struct {
		Columns string   `xml:"COLUMNS"`
		Data    []string `xml:"DATA"`
	}
	xme := XmlMetadataElement{}
	err := d.DecodeElement(&xme, &s)
	if err != nil {
		return &cd, err
	}
	cols := splitCompactRow(xme.Columns, delim)
	for _, rowString := range xme.Data {
		row := splitCompactRow(rowString, delim)
		if len(row) != len(cols) {
			return &cd, errors.New("row length mismatch")
		}
		mapped := make(map[string]string)
		for i, c := range cols {
			mapped[c] = row[i]
		}
		cd.Data = append(cd.Data, mapped)
	}
	return &cd, nil
}

// v needs to be a struct
func ApplyMap(data map[string]string, n interface{}) {
	for k, v := range data {
		f := reflect.ValueOf(n).Elem().FieldByName(k)
		switch f.Kind() {
		case reflect.String:
			x := reflect.ValueOf(n).Elem().FieldByName(k)
			x.Set(reflect.ValueOf(v))
		case reflect.Int:
			parsedInt, _ := strconv.ParseInt(v, 10, 64)
			f.SetInt(parsedInt)
		}
	}
}

func splitCompactRow(row, delim string) []string {
	split := strings.Split(row, delim)
	return split[1 : len(split)-1]
}
