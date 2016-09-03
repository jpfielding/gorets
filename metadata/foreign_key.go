package metadata

import "time"

// MForeignKey ...
type MForeignKey struct {
	Date       time.Time    `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	ForeignKey []ForeignKey `xml:"ForeignKey"`
}

// ForeignKey ...
type ForeignKey struct {
	ForeignKeyID           string `xml:"ForeignKeyID"`
	ParentResourceID       string `xml:"ParentResourceID"`
	ParentClassID          string `xml:"ParentClassID"`
	ParentSystemName       string `xml:"ParentSystemName"`
	ChildResourceID        string `xml:"ChildResourceID"`
	ChildClassID           string `xml:"ChildClassID"`
	ChildSystemName        string `xml:"ChildSystemName"`
	ConditionalParentField string `xml:"ConditionalParentField"`
	ConditionalParentValue string `xml:"ConditionalParentValue"`
	OneToManyFlag          *bool  `xml:"OneToManyFlag"`
}
