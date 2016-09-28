package metadata

import (
	"testing"

	"github.com/jpfielding/gotest/testutils"
)

func TestFieldTransfer(t *testing.T) {
	type Tester struct {
		Bob string
		foo string
		Baz string
	}
	test := &Tester{}
	FieldTransfer(map[string]string{
		"Bob": "foo",
		"Foo": "bar",
		"baz": "blah",
	}).To(test)
	testutils.Equals(t, "foo", test.Bob)
	testutils.Assert(t, "" == test.foo, "shouldnt be able to set lower case fields")
	testutils.Assert(t, "blah" == test.Baz, "should be case insensitive")
}
