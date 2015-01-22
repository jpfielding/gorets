package gorets_metadata

import (
	"time"
)

type MForeignKey struct {
	Version     string
	Date        time.Time
	ForeignKeys []ForeignKey
}

type ForeignKey struct {
	ForeignKeyID           string
	ParentResourceID       string
	ParentClassID          string
	ParentSystemName       string
	ChildResourceID        string
	ChildClassID           string
	ChildSystemName        string
	ConditionalParentField string
	ConditionalParentValue string
	OneToManyFlag          *bool
}
