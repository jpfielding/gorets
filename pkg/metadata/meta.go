package metadata

import (
	"reflect"
	"strings"
)

// MetaInfo provides a meta level for metadata, yeah, sorry
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

// TODO consider adding the common attribute field names along with another
// reflective func to pull all values from a struct into a map

// MISystem ...
var MISystem = MetaInfo{
	Name:        "METADATA-SYSTEM",
	ContentName: "System",
	ContentID:   "SystemID",
	Child:       []MetaInfo{MIResource, MIForeignKey, MIFilter},
}

// MIForeignKey ...
var MIForeignKey = MetaInfo{
	Name:        "METADATA-FOREIGN_KEY",
	ContentName: "ForeignKey",
	ContentID:   "ForeignKeyID",
}

// MIFilter ...
var MIFilter = MetaInfo{
	Name:        "METADATA-FILTER",
	ContentName: "Filter",
	ContentID:   "FilterID",
	Child:       []MetaInfo{MIFilterType},
}

// MIFilterType ...
var MIFilterType = MetaInfo{
	Name:        "METADATA-FILTER_TYPE",
	ContentName: "FilterType",
	ParentID:    "Filter", // FilterID of parent
}

// MIResource ...
var MIResource = MetaInfo{
	Name:        "METADATA-RESOURCE",
	ContentName: "Resource",
	ContentID:   "ResourceID",
	Child: []MetaInfo{
		MIClass,
		MIObject,
		MISearchHelp,
		MIEditMask,
		MIUpdateHelp,
		MILookup,
		MIValidationLookup,
		MIValidationExternal,
		MIValidationExpression,
	},
}

// MIClass ...
var MIClass = MetaInfo{
	Name:        "METADATA-CLASS",
	ContentName: "Class",
	ContentID:   "ClassName",
	Child:       []MetaInfo{MITable, MIUpdate, MIColumnGroup, MIColumnGroupSet},
}

// MITable ...
var MITable = MetaInfo{
	Name:        "METADATA-TABLE",
	ContentName: "Field",
	ContentID:   "MetadataEntryID",
}

// MIObject ...
var MIObject = MetaInfo{
	Name:        "METADATA-OBJECT",
	ContentName: "Object",
	ContentID:   "MetadataEntryID",
}

// MIColumnGroup ...
var MIColumnGroup = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP",
	ContentName: "ColumnGroup",
	ContentID:   "ColumnGroupName",
	Child:       []MetaInfo{MIColumnGroupControl, MIColumnGroupTable, MIColumnGroupNormalization},
}

// MIColumnGroupTable ...
var MIColumnGroupTable = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_TABLE",
	ContentName: "ColumnGroupTable",
	ContentID:   "MetadataEntryID",
}

// MIColumnGroupControl ...
var MIColumnGroupControl = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_CONTROL",
	ContentName: "ColumnGroupControl",
	ContentID:   "MetadataEntryID",
}

// MIColumnGroupNormalization ...
var MIColumnGroupNormalization = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_NORMALIZATION",
	ContentName: "ColumnGroupNormalization",
	ContentID:   "MetadataEntryID",
}

// MIColumnGroupSet ...
var MIColumnGroupSet = MetaInfo{
	Name:        "METADATA-COLUMN_GROUP_SET",
	ContentName: "ColumnGroupSet",
	ContentID:   "ColumnGroupSetName", // MIetadataEntryID also exists ...
}

// MIUpdate ...
var MIUpdate = MetaInfo{
	Name:        "METADATA-UPDATE",
	ContentName: "Update",
	ContentID:   "MetadataEntryID",
	Child:       []MetaInfo{MIUpdateType},
}

// MIUpdateType ...
var MIUpdateType = MetaInfo{
	Name:        "METADATA-UPDATE_TYPE`",
	ContentName: "UpdateType",
	ParentID:    "Update",
}

// MIUpdateHelp ...
var MIUpdateHelp = MetaInfo{
	Name:        "METADATA-UPDATE_HELP`",
	ContentName: "UpdateHelp",
	ContentID:   "UpdateHelpID",
}

// MISearchHelp ...
var MISearchHelp = MetaInfo{
	Name:        "METADATA-SEARCH_HELP",
	ContentName: "SearchHelp",
	ContentID:   "MetadataEntryID",
}

// MIEditMask ...
var MIEditMask = MetaInfo{
	Name:        "METADATA-EDITMASK",
	ContentName: "EditMask",
	ContentID:   "EditMaskID", // MIetadataEntryID also exists ...
}

// MILookup ...
var MILookup = MetaInfo{
	Name:        "METADATA-LOOKUP",
	ContentName: "Lookup",
	ContentID:   "LookupName",
	Child:       []MetaInfo{MILookupType},
}

// MILookupType ...
var MILookupType = MetaInfo{
	Name:        "METADATA-LOOKUP_TYPE",
	ContentName: "LookupType",
	ParentID:    "Lookup",
}

// MIValidationExternal ...
var MIValidationExternal = MetaInfo{
	Name:        "METADATA-VALIDATION_EXTERNAL",
	ContentName: "ValidationExternal",
	ContentID:   "ValidationExternalName", // MIetadataEntryID also exists ...
	Child:       []MetaInfo{MIValidationExternalType},
}

// MIValidationExternalType ...
var MIValidationExternalType = MetaInfo{
	Name:        "METADATA-VALIDATION_EXTERNAL_TYPE",
	ContentName: "ValidationExternalType",
	ContentID:   "MetadataEntryID",
}

// MIValidationExpression ...
var MIValidationExpression = MetaInfo{
	Name:        "METADATA-VALIDATION_EXPRESSION",
	ContentName: "ValidationExpression",
	ContentID:   "ValidationExpressionID", // MIetadataEntryID also exists ...
}

// MIValidationLookup DEPRECATED
var MIValidationLookup = MetaInfo{
	Name:        "METADATA-VALIDATION_LOOKUP",
	ContentName: "ValidationLookup",
	ContentID:   "ValidationLookupName", // MIetadataEntryID also exists ...
	Child:       []MetaInfo{MIValidationLookupType},
}

// MIValidationLookupType DEPRECATED
var MIValidationLookupType = MetaInfo{
	Name:        "METADATA-VALIDATION_LOOKUP_TYPE",
	ContentName: "ValidationLookupType",
	ParentID:    "ValidationLookup",
}
