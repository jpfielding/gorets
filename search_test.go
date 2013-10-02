/**
	provides the searching core
 */
package gorets

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

var compactDecoded string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<COUNT Records="10" />
<DELIMITER value = "09"/>
<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<DATA>	1	2	3	4	5	6	</DATA>
<MAXROWS/>
</RETS>
`

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	cr, err := parseCompactResult(body)
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEqualsInt(t, "bad code", 0, cr.RetsResponse.ReplyCode)
	AssertEquals(t, "bad text", "V2.7.0 2315: Success", cr.RetsResponse.ReplyText)

	AssertEqualsInt(t, "bad count", 10, int(cr.Count))
	AssertEqualsInt(t, "bad header count", 6, len(cr.Columns))

	AssertEquals(t, "bad headers", "A,B,C,D,E,F", strings.Join(cr.Columns,","))

	counter := 0
	for row := range cr.Data {
		if strings.Join(row,",") != "1,2,3,4,5,6" {
			t.Errorf("bad row %s: %s", counter, row)
		}
		counter = counter + 1
	}

	AssertEqualsInt(t, "bad count", 8, counter)
	AssertEqualsTrue(t, "bad max rows", cr.MaxRows)
}

func TestParseStandardXml(t *testing.T) {
}
