package metadata

// MForeignKey ...
type MForeignKey struct {
	Date       DateTime     `xml:"Date,attr,omitempty"`
	Version    Version      `xml:"Version,attr,omitempty"`
	ForeignKey []ForeignKey `xml:"ForeignKey,omitempty"`
}

// ForeignKey ...
type ForeignKey struct {
	ForeignKeyID           RETSID   `xml:"ForeignKeyID,omitempty"`
	ParentResourceID       RETSID   `xml:"ParentResourceID,omitempty"`
	ParentClassID          RETSID   `xml:"ParentClassID,omitempty"`
	ParentSystemName       RETSName `xml:"ParentSystemName,omitempty"`
	ChildResourceID        RETSID   `xml:"ChildResourceID,omitempty"`
	ChildClassID           RETSID   `xml:"ChildClassID,omitempty"`
	ChildSystemName        RETSName `xml:"ChildSystemName,omitempty"`
	ConditionalParentField RETSName `xml:"ConditionalParentField,omitempty"`
	ConditionalParentValue RETSName `xml:"ConditionalParentValue,omitempty"`
	OneToManyFlag          Boolean  `xml:"OneToManyFlag,omitempty"`
}
