package metadata

import "strings"

// Version ...
type Version string

// Major ...
func (mv Version) Major() string {
	return strings.Split(string(mv), ".")[0]
}

// Minor ...
func (mv Version) Minor() string {
	return strings.Split(string(mv), ".")[1]
}

// Release ...
func (mv Version) Release() string {
	return strings.Split(string(mv), ".")[2]
}
