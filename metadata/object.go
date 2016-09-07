package metadata

// MObject ...
type MObject struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`
	Object   []Object `xml:"Object"`
}

// Object ...
type Object struct {
	MetadataEntryID      RETSID            `xml:"MedataEntryID"`
	ObjectType           AlphaNum          `xml:"ObjectType"`
	MIMEType             StringList        `xml:"MIMEType"`
	VisibleName          PlainText         `xml:"VisibleName"`
	Description          PlainText         `xml:"Description"`
	ObjectTimeStamp      RETSName          `xml:"ObjectTimeStamp"`
	ObjectCount          RETSName          `xml:"ObjectCount"`
	LocationAvailability Boolean           `xml:"LocationAvailability"`
	PostSupport          Boolean           `xml:"PostSupport"`
	ObjectData           ResourceClassName `xml:"ObjectData"`
	MaxFileSize          Numeric           `xml:"MaxFileSize"` // TODO limit to positive
}
