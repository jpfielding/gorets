package gorets_metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var compactRsp = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA-SYSTEM Date="Mon, 01 Dec 2014 20:01:22 GMT" Version="39.19.82482">
    <SYSTEM SystemID="RETS" SystemDescription="RETS System" />
    <COMMENTS>RETS SUX</COMMENTS>
  </METADATA-SYSTEM>
</RETS>
`

var expected MSystem = MSystem{
	Version: "39.19.82482",
	Date:    "Mon, 01 Dec 2014 20:01:22 GMT",
	System: System{
		Id:          "RETS",
		Description: "RETS System",
		Comments:    "RETS SUX",
	},
}

// verify it pulls only its section of the common stream
func TestSystemCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactRsp)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-SYSTEM")
	Ok(t, err)

	found := MSystem{}
	err = found.InitFromCompact(parser, start)

	Ok(t, err)

	Equals(t, expected, found)
}

var xmlRsp = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA>
    <METADATA-SYSTEM Date="Mon, 01 Dec 2014 20:01:22 GMT" Version="39.19.82482">
      <System>
        <SystemID>RETS</SystemID>
        <SystemDescription>RETS System</SystemDescription>
        <Comments>RETS SUX</Comments>
      </System>
    </METADATA-SYSTEM>
  </METADATA>
</RETS>
`

// verify it pulls only its section of the common stream
func TestSystemXml(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(xmlRsp)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-SYSTEM")
	Ok(t, err)

	found := MSystem{}
	err = found.InitFromXml(parser, start)

	Ok(t, err)

	Equals(t, expected, found)
}
