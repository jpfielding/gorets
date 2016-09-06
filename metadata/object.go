package metadata

// MObject ...
type MObject struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource string   `xml:"Resource,attr"`
	Object   []Object `xml:"Object"`
}

// Object ...
type Object struct {
	MetadataEntryID      string  `xml:"MedataEntryID"`
	ObjectType           string  `xml:"ObjectType"`
	MimeType             string  `xml:"MIMEType"`
	VisibleName          string  `xml:"VisibleName"`
	Description          string  `xml:"Description"`
	ObjectTimeStamp      string  `xml:"ObjectTimeStamp"`
	ObjectCount          string  `xml:"ObjectCount"`
	LocationAvailability Boolean `xml:"LocationAvailability"`
	PostSupport          Boolean `xml:"PostSupport"`
	ObjectData           string  `xml:"ObjectData"`
	MaxFileSize          Number  `xml:"MaxFileSize"`
}
