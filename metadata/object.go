package metadata

// MObject ...
type MObject struct {
	Date     DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version  Version  `xml:",attr,omitempty" json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty" json:",omitempty"`
	Object   []Object `xml:",omitempty" json:",omitempty"`
}

// Object ...
type Object struct {
	MetadataEntryID      RETSID            `xml:",omitempty" json:",omitempty"`
	ObjectType           AlphaNum          `xml:",omitempty" json:",omitempty"`
	MIMEType             StringList        `xml:",omitempty" json:",omitempty"`
	VisibleName          PlainText         `xml:",omitempty" json:",omitempty"`
	Description          PlainText         `xml:",omitempty" json:",omitempty"`
	ObjectTimeStamp      RETSName          `xml:",omitempty" json:",omitempty"`
	ObjectCount          RETSName          `xml:",omitempty" json:",omitempty"`
	LocationAvailability Boolean           `xml:",omitempty" json:",omitempty"`
	PostSupport          Boolean           `xml:",omitempty" json:",omitempty"`
	ObjectData           ResourceClassName `xml:",omitempty" json:",omitempty"`
	MaxFileSize          Numeric           `xml:",omitempty" json:",omitempty"` // TODO limit to positive
}
