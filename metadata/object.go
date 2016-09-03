package metadata

import "time"

// MObject ...
type MObject struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`
	Object   []Object  `xml:"Object"`
}

// Object ...
type Object struct {
	ObjectType   string `xml:"ObjectType"`
	StandardName string `xml:"ObjectType"`
	MimeType     string `xml:"MimeType"`
	Description  string `xml:"Description"`
}
