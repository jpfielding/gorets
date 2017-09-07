package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "rid", MIResource.ID(test))
	assert.Equal(t, "ln", MILookup.ID(test))
	assert.Equal(t, "cn", MIClass.ID(test))
	assert.Equal(t, "cgn", MIColumnGroup.ID(test))
}

func TestMetaInfoIDMap(t *testing.T) {
	test := map[string]string{
		"ResourceID":      "rid",
		"LookupName":      "ln",
		"MetadataEntryID": "mid",
		"ClassName":       "cn",
		"ColumnGroupName": "cgn",
	}
	assert.Equal(t, "rid", MIResource.ID(test))
	assert.Equal(t, "ln", MILookup.ID(test))
	assert.Equal(t, "cn", MIClass.ID(test))
	assert.Equal(t, "cgn", MIColumnGroup.ID(test))
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
	assert.Equal(t, 25, count(MISystem))
}
