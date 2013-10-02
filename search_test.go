/**
	provides the searching core
 */
package gorets

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	AssertEqualsInt(t, "bad count", 6, len(cr.Columns))

	for row := range cr.Data {
		fmt.Println(row)
	}
}

func TestParseStandardXml(t *testing.T) {
}
