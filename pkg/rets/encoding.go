package rets

import (
	"encoding/xml"
	"io"

	"github.com/jpfielding/gofilters/filter"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// DefaultXMLDecoder the variable used to set a selected charset
var DefaultXMLDecoder = CreateXMLDecoder

// CreateXMLDecoder decodes xml using the given the header if needed
func CreateXMLDecoder(input io.Reader, strict bool) *xml.Decoder {
	// drop any chars that will blow up the xml decoder and replace with a space
	input = filter.NewReader(input, filter.XML10Filter(filter.SpaceChar))
	decoder := xml.NewDecoder(input)
	decoder.Strict = strict
	// this only gets used when a proper xml header is used
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
}

// DefaultReEncodeReader allows overriding the re-encoding operation
var DefaultReEncodeReader = ReEncodeReader

// ReEncodeReader re-encodes a reader based on the http content type provided
func ReEncodeReader(input io.ReadCloser, contentType string) io.ReadCloser {
	if e, _, _ := charset.DetermineEncoding([]byte{}, contentType); e != encoding.Nop {
		type closer struct {
			io.Reader
			io.Closer
		}
		tr := transform.NewReader(input, e.NewDecoder())
		return closer{tr, input}
	}
	return input
}
