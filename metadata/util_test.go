package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "foo", test.Bob)
	assert.Equal(t, "", test.foo, "shouldnt be able to set lower case fields")
	assert.Equal(t, "blah", test.Baz, "should be case insensitive")
}
