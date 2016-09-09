package metadata

// MObject ...
type MObject struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`
	Object   []Object `xml:"Object,omitempty"`
}

// Object ...
type Object struct {
	MetadataEntryID      RETSID            `xml:"MedataEntryID,omitempty"`
	ObjectType           AlphaNum          `xml:"ObjectType,omitempty"`
	MIMEType             StringList        `xml:"MIMEType,omitempty"`
	VisibleName          PlainText         `xml:"VisibleName,omitempty"`
	Description          PlainText         `xml:"Description,omitempty"`
	ObjectTimeStamp      RETSName          `xml:"ObjectTimeStamp,omitempty"`
	ObjectCount          RETSName          `xml:"ObjectCount,omitempty"`
	LocationAvailability Boolean           `xml:"LocationAvailability,omitempty"`
	PostSupport          Boolean           `xml:"PostSupport,omitempty"`
	ObjectData           ResourceClassName `xml:"ObjectData,omitempty"`
	MaxFileSize          Numeric           `xml:"MaxFileSize,omitempty"` // TODO limit to positive
}
