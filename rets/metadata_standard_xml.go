package rets

import (
	"encoding/xml"
	"time"
)

// Represents a <data> element
type StandardXML struct {
	XMLName   xml.Name `xml:"RETS"`
	ReplyCode int      `xml:"ReplyCode"`
	ReplyText string   `xml:"ReplyText"`
	// Entry   []Entry  `xml:"entry"`
}

// Represents an <entry> element
type Entry struct {
	Name     string    `xml:"name"`
	Age      int       `xml:"age"`
	Modified time.Time `xml:"modified,attr"`
}
