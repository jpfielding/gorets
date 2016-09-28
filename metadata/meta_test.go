package metadata

import (
	"testing"

	"github.com/jpfielding/gotest/testutils"
)

func TestMetaInfoIDStruct(t *testing.T) {
	type TestStruct struct {
		ResourceID      string
		MetadataEntryID string
		ClassName       string
		ColumnGroupName string
	}
	test := &TestStruct{
		ResourceID:      "rid",
		MetadataEntryID: "mid",
		ClassName:       "cn",
		ColumnGroupName: "cgn",
	}
	testutils.Equals(t, "rid", MetaResource.ID(test))
	testutils.Equals(t, "mid", MetaLookup.ID(test))
	testutils.Equals(t, "cn", MetaClass.ID(test))
	testutils.Equals(t, "cgn", MetaColumnGroup.ID(test))
}
func TestMetaInfoIDMap(t *testing.T) {
	test := map[string]string{
		"ResourceID":      "rid",
		"MetadataEntryID": "mid",
		"ClassName":       "cn",
		"ColumnGroupName": "cgn",
	}
	testutils.Equals(t, "rid", MetaResource.ID(test))
	testutils.Equals(t, "mid", MetaLookup.ID(test))
	testutils.Equals(t, "cn", MetaClass.ID(test))
	testutils.Equals(t, "cgn", MetaColumnGroup.ID(test))
}
