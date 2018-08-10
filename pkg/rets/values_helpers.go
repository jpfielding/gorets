package rets

import (
	"fmt"
	"net/url"
)

// OptionalStringValue ...
func OptionalStringValue(values url.Values) func(string, string) {
	return func(name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}
}

// OptionalIntValue ...
func OptionalIntValue(values url.Values) func(string, int) {
	return func(name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d", value))
		}
	}
}
