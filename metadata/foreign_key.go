package metadata

// MForeignKey ...
type MForeignKey struct {
	Date       DateTime     `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	ForeignKey []ForeignKey `xml:"ForeignKey"`
}

// ForeignKey ...
type ForeignKey struct {
	ForeignKeyID           RETSID   `xml:"ForeignKeyID"`
	ParentResourceID       RETSID   `xml:"ParentResourceID"`
	ParentClassID          RETSID   `xml:"ParentClassID"`
	ParentSystemName       RETSName `xml:"ParentSystemName"`
	ChildResourceID        RETSID   `xml:"ChildResourceID"`
	ChildClassID           RETSID   `xml:"ChildClassID"`
	ChildSystemName        RETSName `xml:"ChildSystemName"`
	ConditionalParentField RETSName `xml:"ConditionalParentField"`
	ConditionalParentValue RETSName `xml:"ConditionalParentValue"`
	OneToManyFlag          Boolean  `xml:"OneToManyFlag"`
}
