package metadata

// MResource ...
type MResource struct {
	Date     DateTime   `xml:"Date,attr"`
	Version  Version    `xml:"Version,attr"`
	Resource []Resource `xml:"Resource"`
}

// Resource ...
type Resource struct {
	ResourceID                  RETSID    `xml:"ResourceID"`
	StandardName                AlphaNum  `xml:"StandardName"`
	VisibleName                 PlainText `xml:"VisibleName"`
	Description                 PlainText `xml:"Description"`
	KeyField                    RETSName  `xml:"KeyField"`
	ClassCount                  Numeric   `xml:"ClassCount"`
	ClassVersion                Version   `xml:"ClassVersion"`
	ClassDate                   DateTime  `xml:"ClassDate"`
	ObjectVersion               Version   `xml:"ObjectVersion"`
	ObjectDate                  DateTime  `xml:"ObjectDate"`
	SearchHelpVersion           Version   `xml:"SearchHelpVersion"`
	SearchHelpDate              DateTime  `xml:"SearchHelpDate"`
	EditMaskVersion             Version   `xml:"EditMaskVersion"`
	EditMaskDate                DateTime  `xml:"EditMaskDate"`
	LookupVersion               Version   `xml:"LookupVersion"`
	LookupDate                  DateTime  `xml:"LookupDate"`
	UpdateHelpVersion           Version   `xml:"UpdateHelpVersion"`
	UpdateHelpDate              DateTime  `xml:"UpdateHelpDate"`
	ValidationExpressionVersion Version   `xml:"ValidationExpressionVersion"`
	ValidationExpressionDate    DateTime  `xml:"ValidationExpressionDate"`
	ValidationLookupVersion     Version   `xml:"ValidationLookupVersion"`
	ValidationLookupDate        DateTime  `xml:"ValidationLookupDate"`
	ValidationExternalVersion   Version   `xml:"ValidationExternalVersion"`
	ValidationExternalDate      DateTime  `xml:"ValidationExternalDate"`

	// the resource children
	MClass                MClass                `xml:"METADATA-CLASS"`
	MObject               MObject               `xml:"METADATA-OBJECT"`
	MLookup               MLookup               `xml:"METADATA-LOOKUP"`
	MSearchHelp           MSearchHelp           `xml:"METADATA-SEARCH_HELP"`
	MEditMask             MEditMask             `xml:"METADATA-EDIT_MASK"`
	MUpdateHelp           MUpdateHelp           `xml:"METADATA-UPDATE_HELP"`
	MValidationExternal   MValidationExternal   `xml:"METADATA-VALIDATION_EXTERNAL"`
	MValidationExpression MValidationExpression `xml:"METADATA-VALIDATION_EXPRESSION"`

	// deprecated
	MValidationLookup MValidationLookup `xml:"METADATA-VALIDATION_LOOKUP"`
}
