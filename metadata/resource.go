package metadata

// MResource ...
type MResource struct {
	Date     DateTime   `xml:",attr,omitempty" json:",omitempty"`
	Version  Version    `xml:",attr,omitempty" json:",omitempty"`
	Resource []Resource `xml:",omitempty" json:",omitempty"`
}

// Resource ...
type Resource struct {
	ResourceID                  RETSID    `xml:",omitempty" json:",omitempty"`
	StandardName                AlphaNum  `xml:",omitempty" json:",omitempty"`
	VisibleName                 PlainText `xml:",omitempty" json:",omitempty"`
	Description                 PlainText `xml:",omitempty" json:",omitempty"`
	KeyField                    RETSName  `xml:",omitempty" json:",omitempty"`
	ClassCount                  Numeric   `xml:",omitempty" json:",omitempty"`
	ClassVersion                Version   `xml:",omitempty" json:",omitempty"`
	ClassDate                   DateTime  `xml:",omitempty" json:",omitempty"`
	ObjectVersion               Version   `xml:",omitempty" json:",omitempty"`
	ObjectDate                  DateTime  `xml:",omitempty" json:",omitempty"`
	SearchHelpVersion           Version   `xml:",omitempty" json:",omitempty"`
	SearchHelpDate              DateTime  `xml:",omitempty" json:",omitempty"`
	EditMaskVersion             Version   `xml:",omitempty" json:",omitempty"`
	EditMaskDate                DateTime  `xml:",omitempty" json:",omitempty"`
	LookupVersion               Version   `xml:",omitempty" json:",omitempty"`
	LookupDate                  DateTime  `xml:",omitempty" json:",omitempty"`
	UpdateHelpVersion           Version   `xml:",omitempty" json:",omitempty"`
	UpdateHelpDate              DateTime  `xml:",omitempty" json:",omitempty"`
	ValidationExpressionVersion Version   `xml:",omitempty" json:",omitempty"`
	ValidationExpressionDate    DateTime  `xml:",omitempty" json:",omitempty"`
	ValidationLookupVersion     Version   `xml:",omitempty" json:",omitempty"`
	ValidationLookupDate        DateTime  `xml:",omitempty" json:",omitempty"`
	ValidationExternalVersion   Version   `xml:",omitempty" json:",omitempty"`
	ValidationExternalDate      DateTime  `xml:",omitempty" json:",omitempty"`

	// the resource children
	MClass                *MClass                `xml:"METADATA-CLASS,omitempty" json:"METADATA-CLASS,omitempty"`
	MObject               *MObject               `xml:"METADATA-OBJECT,omitempty" json:"METADATA-OBJECT,omitempty"`
	MLookup               *MLookup               `xml:"METADATA-LOOKUP,omitempty" json:"METADATA-LOOKUP,omitempty"`
	MSearchHelp           *MSearchHelp           `xml:"METADATA-SEARCH_HELP,omitempty" json:"METADATA-SEARCH_HELP,omitempty"`
	MEditMask             *MEditMask             `xml:"METADATA-EDITMASK,omitempty" json:"METADATA-EDIT_MASK,omitempty"`
	MUpdateHelp           *MUpdateHelp           `xml:"METADATA-UPDATE_HELP,omitempty" json:"METADATA-UPDATE,omitempty"`
	MValidationExternal   *MValidationExternal   `xml:"METADATA-VALIDATION_EXTERNAL,omitempty" json:"METADATA-VALIDATION_EXTERNAL,omitempty"`
	MValidationExpression *MValidationExpression `xml:"METADATA-VALIDATION_EXPRESSION,omitempty" json:"METADATA-VALIDATION_EXPRESSION,omitempty"`

	// deprecated
	MValidationLookup *MValidationLookup `xml:"METADATA-VALIDATION_LOOKUP,omitempty" json:"METADATA-VALIDAITON_LOOKUP,omitempty"`
}
