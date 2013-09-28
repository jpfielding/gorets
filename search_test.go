/**
	provides the searching core
 */
package gorets

import (
"testing"
)

var compactDecoded string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<COUNT Records="8" />
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

func TestParseCompactDecoded(t *testing.T) {
}
