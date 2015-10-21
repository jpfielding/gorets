package client

import (
	"encoding/xml"
	"io"
)

// SelectedCharsetReader the variable used to set a selected charset
var SelectedCharsetReader func(string, io.Reader) (io.Reader, error)

// GetXMLReader ...
func GetXMLReader(input io.Reader, strict bool) *xml.Decoder {
	decoder := xml.NewDecoder(input)
	if SelectedCharsetReader != nil {
		decoder.CharsetReader = SelectedCharsetReader
	}
	decoder.Strict = strict
	return decoder
}
