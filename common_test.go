package gorets_client

import "testing"


func TestCompactRowParsing(t *testing.T) {
	var col string = `	A	B	C	D	E	F	`
	var row string = `	1	2	3	4		6	`
	var delim string = `	`
	headers := ParseCompactRow(col, delim)
	values := ParseCompactRow(row, delim)

	equals(t, 6, int(len(headers)))
	equals(t, 6, int(len(values)))
}
