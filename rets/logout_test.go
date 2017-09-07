package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessResponseBodyFull(t *testing.T) {
	actual, err := processResponseBody(
		ioutil.NopCloser(strings.NewReader(
			`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				ConnectTime=12345
						 				Billing=Im Billing You
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`,
		)))

	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse{0, "Logging out", 12345, "Im Billing You", "Goodbye"}

	assert.Equal(t, expected, actual)
}

func TestProcessResponseBodyNoBilling(t *testing.T) {
	actual, err := processResponseBody(
		ioutil.NopCloser(strings.NewReader(
			`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				ConnectTime=0
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`,
		)))

	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse{ReplyCode: StatusOK, ReplyText: "Logging out", ConnectTime: 0, SignOffMessage: "Goodbye"}

	assert.Equal(t, expected, actual)
}

func TestProcessResponseBodyNoConnectTime(t *testing.T) {
	actual, err := processResponseBody(
		ioutil.NopCloser(strings.NewReader(
			`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				Billing=Im Billing You
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`,
		)))

	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse{ReplyCode: StatusOK, ReplyText: "Logging out", Billing: "Im Billing You", SignOffMessage: "Goodbye"}

	assert.Equal(t, expected, actual)
}
