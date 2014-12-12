package gorets_client

import (
	"encoding/xml"
	"io"
)

var SelectedCharsetReader func(string, io.Reader) (io.Reader, error) = nil

func GetXmlReader(input io.Reader, strict bool) *xml.Decoder {
	decoder := xml.NewDecoder(input)
	if SelectedCharsetReader != nil {
		decoder.CharsetReader = SelectedCharsetReader
	}
	decoder.Strict = strict
	return decoder
}
