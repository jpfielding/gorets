package rets

import (
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestCompactRowParsing(t *testing.T) {
	var col = `	A	B	C	D	E	F	`
	var row = `	1	2	3	4		6	`
	var delim = `	`
	headers := CompactRow(col).Parse(delim)
	values := CompactRow(row).Parse(delim)

	testutils.Equals(t, 6, int(len(headers)))
	testutils.Equals(t, 6, int(len(values)))
}
