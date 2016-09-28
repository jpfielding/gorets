package metadata

import (
	"reflect"
	"strings"
)

// MetaetaInfo provides a meta level for metadata, yeah, sorry
type MetaInfo struct {
	// Name is the outer element name
	Name string
	// SubName sub elements name for Standard XML
	SubName string
	// SubID id field for either standard or compact sub elements
	SubID string
	// ID in the sub elems for its parent elem
	ParentID string
	// Child are the children of this meta element
	Child []MetaInfo
}

// ID returns the id of the given elem for this meta's info
func (mi MetaInfo) ID(sub interface{}) string {
	if msub, ok := sub.(map[string]string); ok {
		return msub[mi.SubID]
	}
	val := reflect.ValueOf(sub).Elem().FieldByNameFunc(func(n string) bool {
		return strings.ToLower(n) == strings.ToLower(mi.SubID)
	})
	if val.IsValid() {
		return val.String()
	}
	return ""
}

// MetaSystem ...
var MetaSystem MetaInfo = MetaInfo{
	Name:    "METADATA-SYSTEM",
	SubName: "System",
	SubID:   "SystemID",
	Child:   []MetaInfo{MetaResource, MetaForeignKey, MetaFilter},
}

// MetaForeignKey ...
var MetaForeignKey MetaInfo = MetaInfo{
	Name:    "METADATA-FOREIGN_KEY",
	SubName: "ForeignKey",
	SubID:   "ForeignKeyID",
}

// MetaFilter ...
var MetaFilter MetaInfo = MetaInfo{
	Name:    "METADATA-FILTER",
	SubName: "Filter",
	SubID:   "FilterID",
	Child:   []MetaInfo{MetaFilterType},
}

// MetaFilterType ...
var MetaFilterType MetaInfo = MetaInfo{
	Name:     "METADATA-FILTER_TYPE",
	SubName:  "FilterType",
	ParentID: "Filter", // FilterID of parent
}

// MetaResource ...
var MetaResource MetaInfo = MetaInfo{
	Name:    "METADATA-RESOURCE",
	SubName: "Resource",
	SubID:   "ResourceID",
	Child: []MetaInfo{
		MetaClass,
		MetaObject,
		MetaSearchHelp,
		MetaEditMask,
		MetaUpdateHelp,
		MetaLookup,
		MetaValidationLookup,
		MetaValidationExternal,
		MetaValidationExpression,
	},
}

// MetaClass ...
var MetaClass MetaInfo = MetaInfo{
	Name:    "METADATA-CLASS",
	SubName: "Class",
	SubID:   "ClassName",
	Child:   []MetaInfo{MetaTable, MetaUpdate, MetaColumnGroup, MetaColumnGroupSet},
}

// MetaTable ...
var MetaTable MetaInfo = MetaInfo{
	Name:    "METADATA-TABLE",
	SubName: "Field",
	SubID:   "MetadataEntryID",
}

// MetaObject ...
var MetaObject MetaInfo = MetaInfo{
	Name:    "METADATA-OBJECT",
	SubName: "Object",
	SubID:   "MetadataEntryID",
}

// MetaColumnGroup ...
var MetaColumnGroup MetaInfo = MetaInfo{
	Name:    "METADATA-COLUMN_GROUP",
	SubName: "ColumnGroup",
	SubID:   "ColumnGroupName",
	Child:   []MetaInfo{MetaColumnGroupControl, MetaColumnGroupTable, MetaColumnGroupNormalization},
}

// MetaColumnGroupTable ...
var MetaColumnGroupTable MetaInfo = MetaInfo{
	Name:    "METADATA-COLUMN_GROUP_TABLE",
	SubName: "ColumnGroupTable",
	SubID:   "MetadataEntryID",
}

// MetaColumnGroupControl ...
var MetaColumnGroupControl MetaInfo = MetaInfo{
	Name:    "METADATA-COLUMN_GROUP_CONTROL",
	SubName: "ColumnGroupControl",
	SubID:   "MetadataEntryID",
}

// MetaColumnGroupNormalization ...
var MetaColumnGroupNormalization MetaInfo = MetaInfo{
	Name:    "METADATA-COLUMN_GROUP_NORMALIZATION",
	SubName: "ColumnGroupNormalization",
	SubID:   "MetadataEntryID",
}

// MetaColumnGroupSet ...
var MetaColumnGroupSet MetaInfo = MetaInfo{
	Name:    "METADATA-COLUMN_GROUP_SET",
	SubName: "ColumnGroupSet",
	SubID:   "ColumnGroupSetName", // MetaetadataEntryID also exists ...
}

// MetaUpdate ...
var MetaUpdate MetaInfo = MetaInfo{
	Name:    "METADATA-UPDATE",
	SubName: "Update",
	SubID:   "MetadataEntryID",
	Child:   []MetaInfo{MetaUpdateType},
}

// MetaUpdateType ...
var MetaUpdateType MetaInfo = MetaInfo{
	Name:     "METADATA-UPDATE_TYPE`",
	SubName:  "UpdateType",
	ParentID: "Update",
}

// MetaUpdateHelp ...
var MetaUpdateHelp MetaInfo = MetaInfo{
	Name:    "METADATA-UPDATE_HELP`",
	SubName: "UpdateHelp",
	SubID:   "UpdateHelpID",
}

// MetaSearchHelp ...
var MetaSearchHelp MetaInfo = MetaInfo{
	Name:    "METADATA-SEARCH_HELP",
	SubName: "SearchHelp",
	SubID:   "MetadataEntryID",
}

// MetaEditMask ...
var MetaEditMask MetaInfo = MetaInfo{
	Name:    "METADATA-EDIT_MASK",
	SubName: "EditMask",
	SubID:   "EditMaskID", // MetaetadataEntryID also exists ...
}

// MetaLookup ...
var MetaLookup MetaInfo = MetaInfo{
	Name:    "METADATA-LOOKUP",
	SubName: "Lookup",
	SubID:   "MetadataEntryID",
	Child:   []MetaInfo{MetaLookupType},
}

// MetaLookupType ...
var MetaLookupType MetaInfo = MetaInfo{
	Name:     "METADATA-LOOKUP_TYPE",
	SubName:  "LookupType",
	ParentID: "Lookup",
}

// MetaValidationExternal ...
var MetaValidationExternal MetaInfo = MetaInfo{
	Name:    "METADATA-VALIDATION_EXTERNAL",
	SubName: "ValidationExternal",
	SubID:   "ValidationExternalName", // MetaetadataEntryID also exists ...
	Child:   []MetaInfo{MetaValidationExternalType},
}

// MetaValidationExternalType ...
var MetaValidationExternalType MetaInfo = MetaInfo{
	Name:    "METADATA-VALIDATION_EXTERNAL_TYPE",
	SubName: "ValidationExternalType",
	SubID:   "MetadataEntryID",
}

// MetaValidationExpression ...
var MetaValidationExpression MetaInfo = MetaInfo{
	Name:    "METADATA-VALIDATION_EXPRESSION",
	SubName: "ValidationExpression",
	SubID:   "ValidationExpressionID", // MetaetadataEntryID also exists ...
}

// MetaValidationLookup DEPRECATED
var MetaValidationLookup MetaInfo = MetaInfo{
	Name:    "METADATA-VALIDATION_LOOKUP",
	SubName: "ValidationLookup",
	SubID:   "ValidationLookupName", // MetaetadataEntryID also exists ...
	Child:   []MetaInfo{MetaValidationLookupType},
}

// MetaValidationLookupType DEPRECATED
var MetaValidationLookupType MetaInfo = MetaInfo{
	Name:     "METADATA-VALIDATION_LOOKUP_TYPE",
	SubName:  "ValidationLookupType",
	ParentID: "ValidationLookup",
}
