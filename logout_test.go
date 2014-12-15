package gorets_client

import (
	"testing"
)

func TestProcessResponseBodyFull(t *testing.T) {
	actual, err := processResponseBody(`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				ConnectTime=12345
						 				Billing=Im Billing You
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`)
	
	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse { 0, "Logging out", 12345, "Im Billing You", "Goodbye" }

	equals(t, expected, actual)
}

func TestProcessResponseBodyNoBilling(t *testing.T) {
	actual, err := processResponseBody(`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				ConnectTime=0
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`)
	
	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse { ReplyCode: 0, ReplyText: "Logging out", ConnectTime: 0, SignOffMessage: "Goodbye" }

	equals(t, expected, actual)
}

func TestProcessResponseBodyNoConnectTime(t *testing.T) {
	actual, err := processResponseBody(`<RETS ReplyCode="0" ReplyText="Logging out">
						 				<RETS-RESPONSE>
						 				Billing=Im Billing You
						 				SignOffMessage=Goodbye
						 				</RETS-RESPONSE>
		                 				</RETS>`)
	
	if err != nil {
		t.Fail()
	}

	expected := &LogoutResponse { ReplyCode: 0, ReplyText: "Logging out", Billing: "Im Billing You", SignOffMessage: "Goodbye" }

	equals(t, expected, actual)
}