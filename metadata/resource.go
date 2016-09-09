package metadata

// MResource ...
type MResource struct {
	Date     DateTime   `xml:"Date,attr,omitempty"`
	Version  Version    `xml:"Version,attr,omitempty"`
	Resource []Resource `xml:"Resource,omitempty"`
}

// Resource ...
type Resource struct {
	ResourceID                  RETSID    `xml:"ResourceID,omitempty"`
	StandardName                AlphaNum  `xml:"StandardName,omitempty"`
	VisibleName                 PlainText `xml:"VisibleName,omitempty"`
	Description                 PlainText `xml:"Description,omitempty"`
	KeyField                    RETSName  `xml:"KeyField,omitempty"`
	ClassCount                  Numeric   `xml:"ClassCount,omitempty"`
	ClassVersion                Version   `xml:"ClassVersion,omitempty"`
	ClassDate                   DateTime  `xml:"ClassDate,omitempty"`
	ObjectVersion               Version   `xml:"ObjectVersion,omitempty"`
	ObjectDate                  DateTime  `xml:"ObjectDate,omitempty"`
	SearchHelpVersion           Version   `xml:"SearchHelpVersion,omitempty"`
	SearchHelpDate              DateTime  `xml:"SearchHelpDate,omitempty"`
	EditMaskVersion             Version   `xml:"EditMaskVersion,omitempty"`
	EditMaskDate                DateTime  `xml:"EditMaskDate,omitempty"`
	LookupVersion               Version   `xml:"LookupVersion,omitempty"`
	LookupDate                  DateTime  `xml:"LookupDate,omitempty"`
	UpdateHelpVersion           Version   `xml:"UpdateHelpVersion,omitempty"`
	UpdateHelpDate              DateTime  `xml:"UpdateHelpDate,omitempty"`
	ValidationExpressionVersion Version   `xml:"ValidationExpressionVersion,omitempty"`
	ValidationExpressionDate    DateTime  `xml:"ValidationExpressionDate,omitempty"`
	ValidationLookupVersion     Version   `xml:"ValidationLookupVersion,omitempty"`
	ValidationLookupDate        DateTime  `xml:"ValidationLookupDate,omitempty"`
	ValidationExternalVersion   Version   `xml:"ValidationExternalVersion,omitempty"`
	ValidationExternalDate      DateTime  `xml:"ValidationExternalDate,omitempty"`

	// the resource children
	MClass                MClass                `xml:"METADATA-CLASS,omitempty"`
	MObject               MObject               `xml:"METADATA-OBJECT,omitempty"`
	MLookup               MLookup               `xml:"METADATA-LOOKUP,omitempty"`
	MSearchHelp           MSearchHelp           `xml:"METADATA-SEARCH_HELP,omitempty"`
	MEditMask             MEditMask             `xml:"METADATA-EDIT_MASK,omitempty"`
	MUpdateHelp           MUpdateHelp           `xml:"METADATA-UPDATE_HELP,omitempty"`
	MValidationExternal   MValidationExternal   `xml:"METADATA-VALIDATION_EXTERNAL,omitempty"`
	MValidationExpression MValidationExpression `xml:"METADATA-VALIDATION_EXPRESSION,omitempty"`

	// deprecated
	MValidationLookup MValidationLookup `xml:"METADATA-VALIDATION_LOOKUP,omitempty"`
}
