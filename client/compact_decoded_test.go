package client

import (
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestCompactRowParsing(t *testing.T) {
	var col string = `	A	B	C	D	E	F	`
	var row string = `	1	2	3	4		6	`
	var delim string = `	`
	headers := ParseCompactRow(col, delim)
	values := ParseCompactRow(row, delim)

	testutils.Equals(t, 6, int(len(headers)))
	testutils.Equals(t, 6, int(len(values)))
}
