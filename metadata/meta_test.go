package metadata

import (
	"testing"

	"github.com/jpfielding/gotest/testutils"
)

func TestMetaInfoIDStruct(t *testing.T) {
	type TestStruct struct {
		ResourceID      string
		LookupName      string
		MetadataEntryID string
		ClassName       string
		ColumnGroupName string
	}
	test := &TestStruct{
		ResourceID:      "rid",
		LookupName:      "ln",
		MetadataEntryID: "mid",
		ClassName:       "cn",
		ColumnGroupName: "cgn",
	}
	testutils.Equals(t, "rid", MIResource.ID(test))
	testutils.Equals(t, "ln", MILookup.ID(test))
	testutils.Equals(t, "cn", MIClass.ID(test))
	testutils.Equals(t, "cgn", MIColumnGroup.ID(test))
}

func TestMetaInfoIDMap(t *testing.T) {
	test := map[string]string{
		"ResourceID":      "rid",
		"LookupName":      "ln",
		"MetadataEntryID": "mid",
		"ClassName":       "cn",
		"ColumnGroupName": "cgn",
	}
	testutils.Equals(t, "rid", MIResource.ID(test))
	testutils.Equals(t, "ln", MILookup.ID(test))
	testutils.Equals(t, "cn", MIClass.ID(test))
	testutils.Equals(t, "cgn", MIColumnGroup.ID(test))
}

func TestSystemHierarchyCount(t *testing.T) {
	count := func(MetaInfo) int { return 0 }
	count = func(mi MetaInfo) int {
		counter := 1
		for _, c := range mi.Child {
			counter = counter + count(c)
		}
		return counter
	}
	testutils.Equals(t, 25, count(MISystem))
}
