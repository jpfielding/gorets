package metadata

// MValidationExpression ...
type MValidationExpression struct {
	Date     DateTime `xml:",attr,omitempty json:",omitempty"`
	Version  Version  `xml:",attr,omitempty json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty json:",omitempty"`

	ValidationExpression []ValidationExpression `xml:",omitempty json:",omitempty"`
}

// ValidationExpression ...
type ValidationExpression struct {
	MetadataEntryID          RETSID   `xml:",omitempty json:",omitempty"`
	ValidationExpressionID   RETSName `xml:",omitempty json:",omitempty"`
	ValidationExpressionType AlphaNum `xml:",omitempty json:",omitempty"`
	Message                  Text     `xml:",omitempty json:",omitempty"`
	IsCaseSensitive          Boolean  `xml:",omitempty json:",omitempty"`
}
