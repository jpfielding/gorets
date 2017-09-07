package metadata

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	var raw = `<?xml version="1.0" encoding="utf-8"?>
    <RETS ReplyCode="0" ReplyText="Operation successful.">
    <METADATA>
    <METADATA-CLASS Version="01.72.11582" Date="2016-03-29T21:50:11" Resource="Agent">
    </METADATA-CLASS>
    <METADATA-CLASS Version="01.72.11583" Date="2016-03-29T21:50:11" Resource="Office">
    </METADATA-CLASS>
    <METADATA-CLASS Version="01.72.11584" Date="2016-03-29T21:50:11" Resource="Listing">
    </METADATA-CLASS>
    </METADATA>
    </RETS>`
	body := ioutil.NopCloser(strings.NewReader(raw))
	defer body.Close()

	extractor := &Extractor{Body: body}
	response, err := extractor.Open()

	assert.Nil(t, err)
	assert.Equal(t, "Operation successful.", response.ReplyText)

	next := func(resource, version, date string) func(*testing.T) {
		return func(tt *testing.T) {
			mclass := &MClass{}
			err = extractor.DecodeNext("METADATA-CLASS", mclass)
			assert.Nil(t, err)
			assert.Equal(tt, resource, string(mclass.Resource))
			assert.Equal(tt, version, string(mclass.Version))
			assert.Equal(tt, date, string(mclass.Date))
			assert.Equal(tt, 0, len(mclass.Class))
		}
	}

	t.Run("agent", next("Agent", "01.72.11582", "2016-03-29T21:50:11"))
	t.Run("offfice", next("Office", "01.72.11583", "2016-03-29T21:50:11"))
	t.Run("listing", next("Listing", "01.72.11584", "2016-03-29T21:50:11"))

	err = extractor.DecodeNext("METADATA-CLASS", &MClass{})
	assert.Equal(t, io.EOF, err)
}
