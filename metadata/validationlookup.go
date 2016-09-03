package metadata

import "time"

// MValidationLookup ...
type MValidationLookup struct {
	Date             time.Time          `xml:"Date,attr"`
	Version          Version            `xml:"Version,attr"`
	Resource         string             `xml:"Resource,attr"`
	ValidationLookup []ValidationLookup `xml:"ValidationLookup"`
}

// ValidationLookup ...
type ValidationLookup struct {
	Date         time.Time `xml:"Date,attr"`
	Version      Version   `xml:"Version,attr"`
	Parent1Field string    `xml:"Parent1Field"`
	Parent2Field string    `xml:"Parent2Field"`
	Value        string    `xml:"Value"`
}
