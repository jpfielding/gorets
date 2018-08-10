package metadata

import (
	"reflect"
	"strings"
)

// FieldTransfer is a helper to move data from a map to a struct
type FieldTransfer map[string]string

// To is the function for moving the fields to the target
func (fields FieldTransfer) To(target interface{}) {
	for k, v := range fields {
		val := reflect.ValueOf(target).Elem().FieldByNameFunc(func(n string) bool {
			return strings.ToLower(n) == strings.ToLower(k)
		})
		if val.IsValid() && val.CanSet() {
			val.SetString(v)
		}
	}
}
