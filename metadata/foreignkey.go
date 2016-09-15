package metadata

// MForeignKey ...
type MForeignKey struct {
	Date       DateTime     `xml:",attr,omitempty" json:",omitempty"`
	Version    Version      `xml:",attr,omitempty" json:",omitempty"`
	ForeignKey []ForeignKey `xml:",omitempty" json:",omitempty"`
}

// ForeignKey ...
type ForeignKey struct {
	ForeignKeyID           RETSID   `xml:",omitempty" json:",omitempty"`
	ParentResourceID       RETSID   `xml:",omitempty" json:",omitempty"`
	ParentClassID          RETSID   `xml:",omitempty" json:",omitempty"`
	ParentSystemName       RETSName `xml:",omitempty" json:",omitempty"`
	ChildResourceID        RETSID   `xml:",omitempty" json:",omitempty"`
	ChildClassID           RETSID   `xml:",omitempty" json:",omitempty"`
	ChildSystemName        RETSName `xml:",omitempty" json:",omitempty"`
	ConditionalParentField RETSName `xml:",omitempty" json:",omitempty"`
	ConditionalParentValue RETSName `xml:",omitempty" json:",omitempty"`
	OneToManyFlag          Boolean  `xml:",omitempty" json:",omitempty"`
}
