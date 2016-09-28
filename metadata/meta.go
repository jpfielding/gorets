package metadata

import (
	"reflect"
	"strings"
)

// MetaetaInfo provides a meta level for metadata, yeah, sorry
type MetaInfo struct {
	// Name is the outer element name
	Name string
	// ContentName sub elements name for Standard XML
	ContentName string
	// ContentID id field for either standard or compact sub elements
	ContentID string
	// ID in the content elems for its parent elem
	ParentID string
	// Child are the children of this meta element
	Child []MetaInfo
}

// ID returns the id of the given elem for this meta's info
func (mi MetaInfo) ID(sub interface{}) string {
	if msub, ok := sub.(map[string]string); ok {
		return msub[mi.ContentID]
	}
	val := reflect.ValueOf(sub).Elem().FieldByNameFunc(func(n string) bool {
		return strings.ToLower(n) == strings.ToLower(mi.ContentID)
	})
	if val.IsValid() {
		return val.String()
	}
	return ""
}

// MetaSystem ...
var MetaSystem MetaInfo = MetaInfo{
	Name:        "METADATA-SYSTEM",
	ContentName: "System",
	ContentID:   "SystemID",
	Child:       []MetaInfo{MetaResource, MetaForeignKey, MetaFilter},
}

// MetaForeignKey ...
var MetaForeignKey MetaInfo = MetaInfo{
	Name:        "METADATA-FOREIGN_KEY",
	ContentName: "ForeignKey",
	ContentID:   "ForeignKeyID",
}

// MetaFilter ...
var MetaFilter MetaInfo = MetaInfo{
	Name:        "METADATA-FILTER",
	ContentName: "Filter",
	ContentID:   "FilterID",
	Child:       []MetaInfo{MetaFilterType},
}

// MetaFilterType ...
var MetaFilterType MetaInfo = MetaInfo{
	Name:        "METADATA-FILTER_TYPE",
	ContentName: "FilterType",
	ParentID:    "Filter", // FilterID of parent
}

// MetaResource ...
var MetaResource MetaInfo = MetaInfo{
	Name:        "METADATA-RESOURCE",
	ContentName: "Resource",
	ContentID:   "ResourceID",
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
	Name:        "METADATA-CLASS",
	ContentName: "Class",
	ContentID:   "ClassName",
	Child:       []MetaInfo{MetaTable, MetaUpdate, MetaColumnGroup, MetaColumnGroupSet},
}

// MetaTable ...
var MetaTable MetaInfo = MetaInfo{
	Name:        "METADATA-TABLE",
	ContentName: "Field",
	ContentID:   "MetadataEntryID",
}

// MetaObject ...
var MetaObject MetaInfo = MetaInfo{
	Name:        "METADATA-OBJECT",
	ContentName: "Object",
	ContentID:   "MetadataEntryID",
}

// MetaColumnGroup ...
var MetaColumnGroup MetaInfo = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP",
	ContentName: "ColumnGroup",
	ContentID:   "ColumnGroupName",
	Child:       []MetaInfo{MetaColumnGroupControl, MetaColumnGroupTable, MetaColumnGroupNormalization},
}

// MetaColumnGroupTable ...
var MetaColumnGroupTable MetaInfo = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_TABLE",
	ContentName: "ColumnGroupTable",
	ContentID:   "MetadataEntryID",
}

// MetaColumnGroupControl ...
var MetaColumnGroupControl MetaInfo = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_CONTROL",
	ContentName: "ColumnGroupControl",
	ContentID:   "MetadataEntryID",
}

// MetaColumnGroupNormalization ...
var MetaColumnGroupNormalization MetaInfo = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_NORMALIZATION",
	ContentName: "ColumnGroupNormalization",
	ContentID:   "MetadataEntryID",
}

// MetaColumnGroupSet ...
var MetaColumnGroupSet MetaInfo = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_SET",
	ContentName: "ColumnGroupSet",
	ContentID:   "ColumnGroupSetName", // MetaetadataEntryID also exists ...
}

// MetaUpdate ...
var MetaUpdate MetaInfo = MetaInfo{
	Name:        "METADATA-UPDATE",
	ContentName: "Update",
	ContentID:   "MetadataEntryID",
	Child:       []MetaInfo{MetaUpdateType},
}

// MetaUpdateType ...
var MetaUpdateType MetaInfo = MetaInfo{
	Name:        "METADATA-UPDATE_TYPE`",
	ContentName: "UpdateType",
	ParentID:    "Update",
}

// MetaUpdateHelp ...
var MetaUpdateHelp MetaInfo = MetaInfo{
	Name:        "METADATA-UPDATE_HELP`",
	ContentName: "UpdateHelp",
	ContentID:   "UpdateHelpID",
}

// MetaSearchHelp ...
var MetaSearchHelp MetaInfo = MetaInfo{
	Name:        "METADATA-SEARCH_HELP",
	ContentName: "SearchHelp",
	ContentID:   "MetadataEntryID",
}

// MetaEditMask ...
var MetaEditMask MetaInfo = MetaInfo{
	Name:        "METADATA-EDIT_MASK",
	ContentName: "EditMask",
	ContentID:   "EditMaskID", // MetaetadataEntryID also exists ...
}

// MetaLookup ...
var MetaLookup MetaInfo = MetaInfo{
	Name:        "METADATA-LOOKUP",
	ContentName: "Lookup",
	ContentID:   "MetadataEntryID",
	Child:       []MetaInfo{MetaLookupType},
}

// MetaLookupType ...
var MetaLookupType MetaInfo = MetaInfo{
	Name:        "METADATA-LOOKUP_TYPE",
	ContentName: "LookupType",
	ParentID:    "Lookup",
}

// MetaValidationExternal ...
var MetaValidationExternal MetaInfo = MetaInfo{
	Name:        "METADATA-VALIDATION_EXTERNAL",
	ContentName: "ValidationExternal",
	ContentID:   "ValidationExternalName", // MetaetadataEntryID also exists ...
	Child:       []MetaInfo{MetaValidationExternalType},
}

// MetaValidationExternalType ...
var MetaValidationExternalType MetaInfo = MetaInfo{
	Name:        "METADATA-VALIDATION_EXTERNAL_TYPE",
	ContentName: "ValidationExternalType",
	ContentID:   "MetadataEntryID",
}

// MetaValidationExpression ...
var MetaValidationExpression MetaInfo = MetaInfo{
	Name:        "METADATA-VALIDATION_EXPRESSION",
	ContentName: "ValidationExpression",
	ContentID:   "ValidationExpressionID", // MetaetadataEntryID also exists ...
}

// MetaValidationLookup DEPRECATED
var MetaValidationLookup MetaInfo = MetaInfo{
	Name:        "METADATA-VALIDATION_LOOKUP",
	ContentName: "ValidationLookup",
	ContentID:   "ValidationLookupName", // MetaetadataEntryID also exists ...
	Child:       []MetaInfo{MetaValidationLookupType},
}

// MetaValidationLookupType DEPRECATED
var MetaValidationLookupType MetaInfo = MetaInfo{
	Name:        "METADATA-VALIDATION_LOOKUP_TYPE",
	ContentName: "ValidationLookupType",
	ParentID:    "ValidationLookup",
}
